DROP TABLE IF EXISTS rikka.answers;
DROP TABLE IF EXISTS rikka.problems;
DROP TABLE IF EXISTS rikka.notices;
DROP TABLE IF EXISTS rikka.user_profiles;
DROP TABLE IF EXISTS rikka.users;
DROP TABLE IF EXISTS rikka.bastions;
DROP TABLE IF EXISTS rikka.user_groups;

create table notices
(
    id         varchar(191) not null primary key,
    created_at datetime(3)  null,
    updated_at datetime(3)  null,
    source_id  longtext     null,
    title      longtext     null,
    body       longtext     null,
    draft      tinyint(1)   null
);

create table user_groups
(
    id                     varchar(191) not null primary key,
    created_at             datetime(3)  null,
    updated_at             datetime(3)  null,
    name                   varchar(191) not null,
    organization           longtext     not null,
    invitation_code_digest longtext     not null,
    is_full_access         tinyint(1)   not null,
    team_id                longtext     null,
    constraint name unique (name)
);

create table bastions
(
    id            varchar(191) not null primary key,
    created_at    datetime(3)  null,
    updated_at    datetime(3)  null,
    user_group_id varchar(191) not null,
    user          longtext     not null,
    password      longtext     not null,
    host          longtext     not null,
    port          bigint       not null,
    constraint fk_user_groups_bastion foreign key (user_group_id) references user_groups (id)
);

create table users
(
    id              varchar(191) not null primary key,
    created_at      datetime(3)  null,
    updated_at      datetime(3)  null,
    name            varchar(191) not null,
    display_name    longtext     not null,
    password_digest longtext     not null,
    user_group_id   varchar(191) null,
    is_read_only    tinyint(1)   not null,
    constraint name unique (name),
    constraint fk_users_user_group foreign key (user_group_id) references user_groups (id)
);

create table problems
(
    id                  varchar(191)    not null
        primary key,
    created_at          datetime(3)     null,
    updated_at          datetime(3)     null,
    code                varchar(191)    null,
    author_id           varchar(191)    null,
    title               longtext        null,
    body                longtext        null,
    type            longtext null,
    correct_answers longtext null,
    point               bigint unsigned null,
    previous_problem_id varchar(191)    null,
    solved_criterion    bigint unsigned null,
    constraint code
        unique (code),
    constraint fk_problems_author
        foreign key (author_id) references users (id),
    constraint fk_problems_previous_problem
        foreign key (previous_problem_id) references problems (id)
);

create table answers
(
    id            varchar(191)    not null
        primary key,
    created_at    datetime(3)     null,
    updated_at    datetime(3)     null,
    point         bigint unsigned null,
    body          longtext        not null,
    user_group_id varchar(191)    not null,
    problem_id    varchar(191)    not null,
    constraint fk_answers_problem
        foreign key (problem_id) references problems (id),
    constraint fk_answers_user_group
        foreign key (user_group_id) references user_groups (id)
);

create index idx_problems_code
    on problems (code);

create table user_profiles
(
    id                varchar(191) not null
        primary key,
    created_at        datetime(3)  null,
    updated_at        datetime(3)  null,
    user_id           varchar(191) null,
    twitter_id        longtext     null,
    github_id         longtext     null,
    facebook_id       longtext     null,
    self_introduction longtext     null,
    constraint fk_users_user_profile
        foreign key (user_id) references users (id)
);

use rikka;

# user_groups
INSERT INTO user_groups (id, created_at, updated_at, name, organization, invitation_code_digest, is_full_access,
                               team_id)
VALUES ('00000000-0000-4000-a000-000000000000', '1990-01-01 12:34:56.000', '1990-01-01 12:34:56.000', 'admin', 'admin',
        '$2a$10$9i0DvM6.2e7xwT2eueb9ReHKT3mCwRiShNJZM8F/4qvUdrmuncWte', true, 'd446f23f-0e26-4717-9cac-90408e6a2a3b');
INSERT INTO user_groups (id, created_at, updated_at, name, organization, invitation_code_digest, is_full_access,
                               team_id)
VALUES ('00000000-0000-4000-a000-000000000001', '1990-01-01 12:34:56.100', '1990-01-01 12:34:56.100', 'team1',
        'user-org1',
        '$2a$10$9i0DvM6.2e7xwT2eueb9ReHKT3mCwRiShNJZM8F/4qvUdrmuncWte', false, 'd446f23f-0e26-4717-9cac-90408e6a2a3b');
INSERT INTO rikka.user_groups (id, created_at, updated_at, name, organization, invitation_code_digest, is_full_access,
                               team_id)
VALUES ('00000000-0000-4000-a000-000000000002', '1990-01-01 12:34:56.200', '1990-01-01 12:34:56.200', 'team2',
        'user-org2',
        '$2a$10$9i0DvM6.2e7xwT2eueb9ReHKT3mCwRiShNJZM8F/4qvUdrmuncWte', false, 'd446f23f-0e26-4717-9cac-90408e6a2a3b');
INSERT INTO user_groups (id, created_at, updated_at, name, organization, invitation_code_digest, is_full_access,
                         team_id)
VALUES ('00000000-0000-4000-a000-000000000003', '1990-01-01 12:34:56.300', '1990-01-01 12:34:56.300', 'team3',
        'user-org3',
        '$2a$10$9i0DvM6.2e7xwT2eueb9ReHKT3mCwRiShNJZM8F/4qvUdrmuncWte', false, 'd446f23f-0e26-4717-9cac-90408e6a2a3b');


# bastions
## user_groups_id: 00000000-0000-4000-a000-000000000000
INSERT INTO rikka.bastions (id, created_at, updated_at, user_group_id, user, password, host, port)
VALUES ('00000000-0000-4000-a000-000000000000', '1990-01-01 12:34:56.000', '1990-01-01 12:34:56.000',
        '00000000-0000-4000-a000-000000000000', 'admin', 'password', 'host', 22);
## user_groups_id: 00000000-0000-4000-a000-000000000001
INSERT INTO rikka.bastions (id, created_at, updated_at, user_group_id, user, password, host, port)
VALUES ('00000000-0000-4000-a000-000000000001', '1990-01-01 12:34:56.100', '1990-01-01 12:34:56.100',
        '00000000-0000-4000-a000-000000000001', 'user1', 'password', 'host', 22);
## user_groups_id: 00000000-0000-4000-a000-000000000002
INSERT INTO rikka.bastions (id, created_at, updated_at, user_group_id, user, password, host, port)
VALUES ('00000000-0000-4000-a000-000000000002', '1990-01-01 12:34:56.200', '1990-01-01 12:34:56.200',
        '00000000-0000-4000-a000-000000000002', 'user2', 'password', 'host', 22);
## user_groups_id: 00000000-0000-4000-a000-000000000003
INSERT INTO rikka.bastions (id, created_at, updated_at, user_group_id, user, password, host, port)
VALUES ('00000000-0000-4000-a000-000000000003', '1990-01-01 12:34:56.300', '1990-01-01 12:34:56.300',
        '00000000-0000-4000-a000-000000000003', 'user3', 'password', 'host', 22);


# users
## user_groups_id: 00000000-0000-4000-a000-000000000000
INSERT INTO rikka.users (id, created_at, updated_at, name, display_name, password_digest, user_group_id, is_read_only)
VALUES ('00000000-0000-4000-a000-000000000000', '1990-01-01 12:34:56.000', '1990-01-01 12:34:56.000', 'admin', 'admin',
        '$2a$10$2h4slyxP75YpEmndafD.GOt/hzhRWs.lvLWGIq0/CnGgSSpQnQBAq',
        '00000000-0000-4000-a000-000000000000', false);
## user_groups_id: 00000000-0000-4000-a000-000000000001
INSERT INTO rikka.users (id, created_at, updated_at, name, display_name, password_digest, user_group_id, is_read_only)
VALUES ('00000000-0000-4000-a000-000000000001', '1990-01-01 12:34:56.100', '1990-01-01 12:34:56.100', 'user1', 'user1',
        '$2a$10$2h4slyxP75YpEmndafD.GOt/hzhRWs.lvLWGIq0/CnGgSSpQnQBAq',
        '00000000-0000-4000-a000-000000000001', false);
INSERT INTO rikka.users (id, created_at, updated_at, name, display_name, password_digest, user_group_id, is_read_only)
VALUES ('00000000-0000-4000-a000-000000000002', '1990-01-01 12:34:56.200', '1990-01-01 12:34:56.200', 'user2', 'user2',
        '$2a$10$2h4slyxP75YpEmndafD.GOt/hzhRWs.lvLWGIq0/CnGgSSpQnQBAq',
        '00000000-0000-4000-a000-000000000001', false);
INSERT INTO rikka.users (id, created_at, updated_at, name, display_name, password_digest, user_group_id, is_read_only)
VALUES ('00000000-0000-4000-a000-000000000003', '1990-01-01 12:34:56.300', '1990-01-01 12:34:56.300', 'user3', 'user3',
        '$2a$10$2h4slyxP75YpEmndafD.GOt/hzhRWs.lvLWGIq0/CnGgSSpQnQBAq',
        '00000000-0000-4000-a000-000000000001', true);
## user_groups_id: 00000000-0000-4000-a000-000000000002
INSERT INTO rikka.users (id, created_at, updated_at, name, display_name, password_digest, user_group_id, is_read_only)
VALUES ('00000000-0000-4000-a000-000000000004', '1990-01-01 12:34:56.400', '1990-01-01 12:34:56.400', 'user4', 'user4',
        '$2a$10$2h4slyxP75YpEmndafD.GOt/hzhRWs.lvLWGIq0/CnGgSSpQnQBAq',
        '00000000-0000-4000-a000-000000000002', false);
INSERT INTO rikka.users (id, created_at, updated_at, name, display_name, password_digest, user_group_id, is_read_only)
VALUES ('00000000-0000-4000-a000-000000000005', '1990-01-01 12:34:56.500', '1990-01-01 12:34:56.500', 'user5', 'user5',
        '$2a$10$2h4slyxP75YpEmndafD.GOt/hzhRWs.lvLWGIq0/CnGgSSpQnQBAq',
        '00000000-0000-4000-a000-000000000002', false);
INSERT INTO rikka.users (id, created_at, updated_at, name, display_name, password_digest, user_group_id, is_read_only)
VALUES ('00000000-0000-4000-a000-000000000006', '1990-01-01 12:34:56.600', '1990-01-01 12:34:56.600', 'user6', 'user6',
        '$2a$10$2h4slyxP75YpEmndafD.GOt/hzhRWs.lvLWGIq0/CnGgSSpQnQBAq',
        '00000000-0000-4000-a000-000000000002', true);
## user_groups_id: 00000000-0000-4000-a000-000000000003
INSERT INTO rikka.users (id, created_at, updated_at, name, display_name, password_digest, user_group_id, is_read_only)
VALUES ('00000000-0000-4000-a000-000000000007', '1990-01-01 12:34:56.700', '1990-01-01 12:34:56.700', 'user7', 'user7',
        '$2a$10$2h4slyxP75YpEmndafD.GOt/hzhRWs.lvLWGIq0/CnGgSSpQnQBAq',
        '00000000-0000-4000-a000-000000000003', false);
INSERT INTO rikka.users (id, created_at, updated_at, name, display_name, password_digest, user_group_id, is_read_only)
VALUES ('00000000-0000-4000-a000-000000000008', '1990-01-01 12:34:56.800', '1990-01-01 12:34:56.800', 'user8', 'user8',
        '$2a$10$2h4slyxP75YpEmndafD.GOt/hzhRWs.lvLWGIq0/CnGgSSpQnQBAq',
        '00000000-0000-4000-a000-000000000003', false);
INSERT INTO rikka.users (id, created_at, updated_at, name, display_name, password_digest, user_group_id, is_read_only)
VALUES ('00000000-0000-4000-a000-000000000009', '1990-01-01 12:34:56.900', '1990-01-01 12:34:56.900', 'user9', 'user9',
        '$2a$10$2h4slyxP75YpEmndafD.GOt/hzhRWs.lvLWGIq0/CnGgSSpQnQBAq',
        '00000000-0000-4000-a000-000000000003', true);

# user_profiles
## user_id: 00000000-0000-4000-a000-000000000001
INSERT INTO rikka.user_profiles (id, created_at, updated_at, user_id, twitter_id, github_id, facebook_id,
                                 self_introduction)
VALUES ('00000000-0000-4000-a000-000000000000', '1990-01-01 12:34:56.000', '1990-01-01 12:34:56.000',
        '00000000-0000-4000-a000-000000000001', 'tid', 'gid', 'fid', '自己紹介内容1');
## user_id: 00000000-0000-4000-a000-000000000004
INSERT INTO rikka.user_profiles (id, created_at, updated_at, user_id, twitter_id, github_id, facebook_id,
                                 self_introduction)
VALUES ('00000000-0000-4000-a000-000000000001', '1990-01-01 12:34:56.100', '1990-01-01 12:34:56.100',
        '00000000-0000-4000-a000-000000000004', 'tid', 'gid', 'fid', '自己紹介内容2');
## user_id: 00000000-0000-4000-a000-000000000007
INSERT INTO rikka.user_profiles (id, created_at, updated_at, user_id, twitter_id, github_id, facebook_id,
                                 self_introduction)
VALUES ('00000000-0000-4000-a000-000000000002', '1990-01-01 12:34:56.200', '1990-01-01 12:34:56.200',
        '00000000-0000-4000-a000-000000000007', 'tid', 'gid', 'fid', '自己紹介内容3');


# notices
INSERT INTO rikka.notices (id, created_at, updated_at, source_id, title, body, draft)
VALUES ('00000000-0000-4000-a000-000000000000', '1990-01-01 12:34:56.000', '1990-01-01 12:34:56.000',
        '00000000-0000-4000-a000-000000000000', '通知タイトル1', '通知メッセージ1', false);
INSERT INTO rikka.notices (id, created_at, updated_at, source_id, title, body, draft)
VALUES ('00000000-0000-4000-a000-000000000001', '1990-01-01 12:34:56.100', '1990-01-01 12:34:56.100',
        '00000000-0000-4000-a000-000000000001', '通知タイトル2', '通知メッセージ2', false);
INSERT INTO rikka.notices (id, created_at, updated_at, source_id, title, body, draft)
VALUES ('00000000-0000-4000-a000-000000000002', '1990-01-01 12:34:56.200', '1990-01-01 12:34:56.200',
        '00000000-0000-4000-a000-000000000002', '通知タイトル3', '通知メッセージ3', false);


# problems
INSERT INTO rikka.problems (id, created_at, updated_at, code, author_id, title, body, point, previous_problem_id,
                            solved_criterion)
VALUES ('00000000-0000-4000-a000-000000000000', '1990-01-01 12:34:56.000', '1990-01-01 12:34:56.000', 'abc',
        '00000000-0000-4000-a000-000000000000', '問題タイトル1', '---
code: abc
title: テスト問題
point: 100
solvedCriterion: 100
connectInfo:
  - hostname: post
    command: 192.168.100.1
    user: user
    password: password
    port: 22
    type: ssh
---
問題内容1', 100, null, 100);
INSERT INTO rikka.problems (id, created_at, updated_at, code, author_id, title, body, point, previous_problem_id,
                            solved_criterion)
VALUES ('00000000-0000-4000-a000-000000000001', '1990-01-01 12:34:56.100', '1990-01-01 12:34:56.100', 'def',
        '00000000-0000-4000-a000-000000000000', '問題タイトル2', '問題内容2', 200, null, 200);
INSERT INTO rikka.problems (id, created_at, updated_at, code, author_id, title, body, point, previous_problem_id,
                            solved_criterion)
VALUES ('00000000-0000-4000-a000-000000000002', '1990-01-01 12:34:56.200', '1990-01-01 12:34:56.200', 'ghi',
        '00000000-0000-4000-a000-000000000000', '問題タイトル3', '問題内容3', 300, null, 300);


# answers
INSERT INTO rikka.answers (id, created_at, updated_at, point, body, user_group_id, problem_id)
VALUES ('00000000-0000-4000-a000-000000000000', '1990-01-01 12:34:56.000', '1990-01-01 12:34:56.000', 0, '回答内容1',
        '00000000-0000-4000-a000-000000000001', '00000000-0000-4000-a000-000000000000');
INSERT INTO rikka.answers (id, created_at, updated_at, point, body, user_group_id, problem_id)
VALUES ('00000000-0000-4000-a000-000000000001', '1990-01-01 12:34:56.100', '1990-01-01 12:34:56.100', null, '回答内容2',
        '00000000-0000-4000-a000-000000000001', '00000000-0000-4000-a000-000000000000');
INSERT INTO rikka.answers (id, created_at, updated_at, point, body, user_group_id, problem_id)
VALUES ('00000000-0000-4000-a000-000000000002', '1990-01-01 12:34:56.200', '1990-01-01 12:34:56.200', 100, '回答内容3',
        '00000000-0000-4000-a000-000000000001', '00000000-0000-4000-a000-000000000000');
INSERT INTO rikka.answers (id, created_at, updated_at, point, body, user_group_id, problem_id)
VALUES ('00000000-0000-4000-a000-000000000003', '1990-01-01 12:34:56.300', '1990-01-01 12:34:56.300', 0, '回答内容4',
        '00000000-0000-4000-a000-000000000001', '00000000-0000-4000-a000-000000000001');
INSERT INTO rikka.answers (id, created_at, updated_at, point, body, user_group_id, problem_id)
VALUES ('00000000-0000-4000-a000-000000000004', '1990-01-01 12:34:56.400', '1990-01-01 12:34:56.400', null, '回答内容5',
        '00000000-0000-4000-a000-000000000001', '00000000-0000-4000-a000-000000000001');
INSERT INTO rikka.answers (id, created_at, updated_at, point, body, user_group_id, problem_id)
VALUES ('00000000-0000-4000-a000-000000000005', '1990-01-01 12:34:56.500', '1990-01-01 12:34:56.500', 200, '回答内容6',
        '00000000-0000-4000-a000-000000000001', '00000000-0000-4000-a000-000000000001');
INSERT INTO rikka.answers (id, created_at, updated_at, point, body, user_group_id, problem_id)
VALUES ('00000000-0000-4000-a000-000000000006', '1990-01-01 12:34:56.600', '1990-01-01 12:34:56.600', 0, '回答内容7',
        '00000000-0000-4000-a000-000000000001', '00000000-0000-4000-a000-000000000002');
INSERT INTO rikka.answers (id, created_at, updated_at, point, body, user_group_id, problem_id)
VALUES ('00000000-0000-4000-a000-000000000007', '1990-01-01 12:34:56.700', '1990-01-01 12:34:56.700', null, '回答内容8',
        '00000000-0000-4000-a000-000000000001', '00000000-0000-4000-a000-000000000002');
INSERT INTO rikka.answers (id, created_at, updated_at, point, body, user_group_id, problem_id)
VALUES ('00000000-0000-4000-a000-000000000008', '1990-01-01 12:34:56.800', '1990-01-01 12:34:56.800', 300, '回答内容9',
        '00000000-0000-4000-a000-000000000001', '00000000-0000-4000-a000-000000000002');