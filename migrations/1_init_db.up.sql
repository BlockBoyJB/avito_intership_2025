create table if not exists teams
(
    name varchar primary key
);

create table if not exists users
(
    id        varchar primary key,
    username  varchar                                           not null,
    team_name varchar references teams (name) on delete cascade not null,
    is_active bool default true                                 not null
);

create table if not exists pull_requests
(
    id         varchar primary key,
    name       varchar                                                       not null,
    author_id  varchar references users (id) on delete cascade               not null,
    status     varchar check ( status in ('OPEN', 'MERGED') ) default 'OPEN' not null,
    created_at timestamptz                                    default now()  not null,
    merged_at  timestamptz                                                   null
);

create table if not exists pr_reviewers
(
    id          serial primary key,
    pr_id       varchar references pull_requests (id) on delete cascade not null,
    reviewer_id varchar references users (id) on delete cascade         not null,
    unique (pr_id, reviewer_id)
)