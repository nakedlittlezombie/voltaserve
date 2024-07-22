use sea_orm_migration::prelude::*;

use crate::models::v1::{
    Group, Grouppermission,
};

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(
        &self,
        manager: &SchemaManager,
    ) -> Result<(), DbErr> {
        manager
            .create_table(
                Table::create()
                    .table(Grouppermission::Table)
                    .col(
                        ColumnDef::new(Grouppermission::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(ColumnDef::new(Grouppermission::GroupId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(Grouppermission::Table, Grouppermission::GroupId)
                            .to(Group::Table, Group::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(Grouppermission::ResourceId).text())
                    .col(ColumnDef::new(Grouppermission::Permission).text())
                    .col(
                        ColumnDef::new(Grouppermission::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .table(Grouppermission::Table)
                    .col(Grouppermission::GroupId)
                    .col(Grouppermission::ResourceId)
                    .unique()
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .table(Grouppermission::Table)
                    .col(Grouppermission::GroupId)
                    .to_owned(),
            )
            .await?;
        manager
            .create_index(
                Index::create()
                    .table(Grouppermission::Table)
                    .col(Grouppermission::ResourceId)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }

    async fn down(
        &self,
        manager: &SchemaManager,
    ) -> Result<(), DbErr> {
        manager
            .drop_table(
                Table::drop()
                    .table(Grouppermission::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}

