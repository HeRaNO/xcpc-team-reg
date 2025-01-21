CREATE TABLE t_user (
    user_id      BIGSERIAL PRIMARY KEY,
    user_name    VARCHAR(50) NOT NULL,
    email        VARCHAR(50),
    school       INT,
    stu_id       VARCHAR(15),
    belong_team  BIGINT,
    tshirt       VARCHAR(3),
    is_uestc_stu INT
);

CREATE INDEX idx_user_belong_team ON t_user (belong_team);

CREATE TABLE t_auth (
    user_id BIGINT PRIMARY KEY,
    email   VARCHAR(50),
    pwd     VARCHAR(100)
);

CREATE TABLE t_team (
    team_id          BIGSERIAL PRIMARY KEY,
    team_name        VARCHAR(50) NOT NULL UNIQUE,
    member_cnt       INT,
    team_account     VARCHAR(20),
    team_password    VARCHAR(20),
    team_affiliation VARCHAR(50),
    invite_token     VARCHAR(10)
);
