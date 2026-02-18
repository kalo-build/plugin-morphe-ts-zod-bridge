name: Person
fields:
  ID:
    type: AutoIncrement
  Name:
    type: String
identifiers:
  primary: ID
  name: Name
related:
  Note:
    type: HasOnePoly
    through: Commentable
    aliased: Comment
