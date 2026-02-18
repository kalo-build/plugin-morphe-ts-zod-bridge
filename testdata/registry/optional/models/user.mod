name: User
fields:
  ID:
    type: AutoIncrement
  Email:
    type: String
  Nickname:
    type: String
    attributes:
      - optional
  Bio:
    type: String
    attributes:
      - optional
identifiers:
  primary: ID
  email: Email
related:
  Organization:
    type: ForOne
