name: Organization
fields:
  ID:
    type: AutoIncrement
  Name:
    type: String
identifiers:
  primary: ID
  name: Name
related:
  User:
    type: HasMany
