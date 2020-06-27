package pwapi

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
	"pathwar.land/pathwar/v2/go/pkg/pwsso"
)

func (svc *service) UserGetSession(ctx context.Context, _ *UserGetSession_Input) (*UserGetSession_Output, error) {
	token, err := tokenFromContext(ctx)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}
	svc.logger.Debug("token", zap.Any("token", token))

	output := &UserGetSession_Output{}
	output.Claims = pwsso.ClaimsFromToken(token)

	// try loading it from database
	output.User, err = svc.loadOAuthUser(output.Claims.ActionToken.Sub)
	if err != nil && !pwdb.IsRecordNotFoundError(err) {
		return nil, errcode.ErrGetOAuthUser.Wrap(err)
	}

	// new user
	if pwdb.IsRecordNotFoundError(err) {
		output.IsNewUser = true
		if _, err = svc.newUserFromClaims(output.Claims); err != nil {
			return nil, errcode.ErrNewUserFromClaims.Wrap(err)
		}
		if output.User, err = svc.loadOAuthUser(output.Claims.ActionToken.Sub); err != nil {
			return nil, errcode.ErrGetOAuthUser.Wrap(err)
		}
	}

	if output.User.Username != output.Claims.PreferredUsername {
		return nil, errcode.ErrDifferentUserBetweenTokenAndDatabase
	}

	// FIXME: output.Notifications = COUNT
	output.Notifications = int32(rand.Intn(10))

	output.Seasons, err = svc.loadUserSeasons(ctx)
	if err != nil {
		return nil, errcode.ErrLoadUserSeasons.Wrap(err)
	}

	return output, nil
}

func (svc *service) loadUserSeasons(ctx context.Context) ([]*UserGetSession_Output_SeasonAndTeam, error) {
	var (
		seasons     []*pwdb.Season
		memberships []*pwdb.TeamMember
	)

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated
	}

	// get season organizations for user
	err = svc.db.
		Preload("Team").
		Preload("Team.Organization").
		Where(pwdb.TeamMember{UserID: userID}).
		Find(&memberships).
		Error
	if err != nil {
		return nil, errcode.ErrGetUserOrganizations.Wrap(err)
	}

	// get all available seasons
	err = svc.db.
		Where(pwdb.Season{Visibility: pwdb.Season_Public}).
		// FIXME: admins can see everything
		Find(&seasons).
		Error
	if err != nil {
		return nil, errcode.ErrGetSeasons.Wrap(err)
	}

	output := []*UserGetSession_Output_SeasonAndTeam{}
	for _, season := range seasons {
		item := &UserGetSession_Output_SeasonAndTeam{
			Season: season,
		}

		for _, membership := range memberships {
			if membership.Team.SeasonID == season.ID {
				item.Team = membership.Team
				break
			}
		}

		output = append(output, item)
	}

	return output, nil
}

func (svc *service) loadOAuthUser(subject string) (*pwdb.User, error) {
	var user pwdb.User
	err := svc.db.
		Preload("ActiveTeamMember").
		Preload("ActiveTeamMember.Team").
		Preload("ActiveTeamMember.Team.Season").
		Preload("ActiveTeamMember.Team.Organization").
		Where(pwdb.User{OAuthSubject: subject}).
		First(&user).
		Error
	if err != nil {
		return nil, errcode.ErrGetUserBySubject.Wrap(err)
	}

	return &user, nil
}

func (svc *service) newUserFromClaims(claims *pwsso.Claims) (*pwdb.User, error) {
	if !claims.EmailVerified {
		return nil, errcode.ErrEmailAddressNotVerified
	}

	gravatarURL := fmt.Sprintf("https://www.gravatar.com/avatar/%x", md5.Sum([]byte(claims.Email)))

	var season pwdb.Season
	if err := svc.db.Where(pwdb.Season{IsDefault: true}).First(&season).Error; err != nil {
		return nil, errcode.ErrGetDefaultSeason.Wrap(err)
	}

	user := pwdb.User{
		Username:     claims.PreferredUsername,
		Email:        claims.Email,
		GravatarURL:  gravatarURL,
		OAuthSubject: claims.ActionToken.Sub,
		// WebsiteURL
		// Locale

		TeamMemberships:         []*pwdb.TeamMember{},
		OrganizationMemberships: []*pwdb.OrganizationMember{},
		DeletionStatus:          pwdb.DeletionStatus_Active,
	}
	organization := pwdb.Organization{
		Name:           claims.PreferredUsername,
		GravatarURL:    gravatarURL,
		DeletionStatus: pwdb.DeletionStatus_Active,
		SoloSeason:     true,
		// Locale
	}
	organizationMember := pwdb.OrganizationMember{
		//User: &user,
		Organization: &organization,
		Role:         pwdb.OrganizationMember_Owner,
	}
	team := pwdb.Team{
		Season:         &season,
		IsDefault:      true,
		Organization:   &organization,
		DeletionStatus: pwdb.DeletionStatus_Active,
	}
	seasonMember := pwdb.TeamMember{
		User: &user,
		Team: &team,
		Role: pwdb.TeamMember_Owner,
	}
	user.OrganizationMemberships = []*pwdb.OrganizationMember{&organizationMember}

	err := svc.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&seasonMember).Error; err != nil {
			return err
		}
		activity := pwdb.Activity{
			Kind:           pwdb.Activity_UserRegister,
			AuthorID:       user.ID, // FIXME: author should be based on ctx
			UserID:         user.ID,
			TeamID:         team.ID,
			OrganizationID: organization.ID,
			TeamMemberID:   seasonMember.ID,
			SeasonID:       season.ID,
		}
		if err := tx.Create(&activity).Error; err != nil {
			return err
		}

		// set active season
		err := tx.Model(&user).Updates(pwdb.User{
			ActiveTeamMemberID: seasonMember.ID,
			ActiveSeasonID:     season.ID,
		}).
			Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errcode.ErrCommitUserTransaction.Wrap(err)
	}

	return &user, nil
}
