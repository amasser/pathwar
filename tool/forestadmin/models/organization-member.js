// This model was generated by Lumber. However, you remain in control of your models.
// Learn how here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/models/enrich-your-models
module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  // This section contains the fields of your model, mapped to your table's columns.
  // Learn more here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/models/enrich-your-models#declaring-a-new-field-in-a-model
  const OrganizationMember = sequelize.define('organizationMember', {
    createdAt: {
      type: DataTypes.DATE,
    },
    updatedAt: {
      type: DataTypes.DATE,
    },
    role: {
      type: DataTypes.INTEGER,
    },
  }, {
    tableName: 'organization_member',
    underscored: true,
  });

  // This section contains the relationships for this model. See: https://docs.forestadmin.com/documentation/v/v6/reference-guide/relationships#adding-relationships.
  OrganizationMember.associate = (models) => {
    OrganizationMember.belongsTo(models.user, {
      foreignKey: {
        name: 'userIdKey',
        field: 'user_id',
      },
      as: 'user',
    });
    OrganizationMember.belongsTo(models.organization, {
      foreignKey: {
        name: 'organizationIdKey',
        field: 'organization_id',
      },
      as: 'organization',
    });
  };

  return OrganizationMember;
};
