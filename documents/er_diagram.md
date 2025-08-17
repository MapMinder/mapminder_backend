# Database Schema

```mermaid
---
config:
    layout: elk
---
erDiagram

oauth_tokens {
    int oauth_id PK
    varchar(40) user_id FK 
    varchar(50) oauth_provider
    varchar(255) oauth_provider_id
    text refresh_token
    timestamp access_token_expiry
    datetime created_at
    datetime updated_at
}

user_setting {
    varchar(40) user_setting_id PK
    varchar(40) user_id FK
    varchar(10) theme
    varchar(3) language
    datetime created_at
    datetime updated_at
}

users {
    varchar(40) user_id PK
    varchar(60) username
    varchar(255) email
    datetime created_at
    datetime updated_at
}

tags {
    int tag_id PK
    id note_id FK
    varchar(20) tag
}

locations {
    int location_id PK
    int reminder_id FK
    float latitude
    float longitude
    varchar(60) place_name
}

reminders {
    int reminder_id PK
    varchar(40) user_id FK
    int note_id
    int location_id
    bool reminder_flag
    datetime date_time
}

notes {
    int note_id PK
    int reminder_id FK
    varchar(255) title
    text content
    datetime created_at
    datetime updated_at
}

notifications {
    int notification_id PK
    varchar(40) user_id FK
    int reminder_id FK
    datetime sent_at
    varchar(10) status
}

reminder_recurrences {
    int recurrence_id PK
    int reminder_id FK
    varchar(10) recurrence_type
}

users ||--o{ reminders : has
users ||--|| oauth_tokens : has
users ||--|| user_setting : has
reminders ||--o{ notes : has
notes ||--o{ tags : has
reminders ||--|| locations : has
reminders ||--|| notifications : has
reminders ||--o| reminder_recurrences : might_have
```
