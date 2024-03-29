-- +goose Up

CREATE SCHEMA donoengine;

-- Language

CREATE TABLE donoengine.language (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    ietf            text                        NOT NULL,
    name            text                        NOT NULL
);

ALTER TABLE ONLY donoengine.language ADD CONSTRAINT language_ietf_key
    UNIQUE (ietf);

ALTER TABLE ONLY donoengine.language ADD CONSTRAINT language_ietf_check
    CHECK (ietf <> '' AND length(ietf) <= 12);
ALTER TABLE ONLY donoengine.language ADD CONSTRAINT language_name_check
    CHECK (name <> '' AND length(name) <= 24);

-- Website

CREATE TABLE donoengine.website (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    domain          text                        NOT NULL,
    name            text                        NOT NULL
);

ALTER TABLE ONLY donoengine.website ADD CONSTRAINT website_domain_key
    UNIQUE (domain);

ALTER TABLE ONLY donoengine.website ADD CONSTRAINT website_domain_check
    CHECK (domain <> '' AND length(domain) <= 32);
ALTER TABLE ONLY donoengine.website ADD CONSTRAINT website_name_check
    CHECK (name <> '' AND length(name) <= 48);

-- Category Type

CREATE TABLE donoengine.category_type (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    code            text                        NOT NULL,
    name            text                        NOT NULL
);

ALTER TABLE ONLY donoengine.category_type ADD CONSTRAINT category_type_code_key
    UNIQUE (code);

ALTER TABLE ONLY donoengine.category_type ADD CONSTRAINT category_type_code_check
    CHECK (code <> '' AND length(code) <= 24);
ALTER TABLE ONLY donoengine.category_type ADD CONSTRAINT category_type_name_check
    CHECK (name <> '' AND length(name) <= 24);

-- Category

CREATE TABLE donoengine.category (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    type_id         bigint                      NOT NULL,
    code            text                        NOT NULL,
    name            text                        NOT NULL
);

ALTER TABLE ONLY donoengine.category ADD CONSTRAINT category_type_id_fkey
    FOREIGN KEY (type_id) REFERENCES donoengine.category_type(id);

ALTER TABLE ONLY donoengine.category ADD CONSTRAINT category_type_id_code_key
    UNIQUE (type_id, code);

ALTER TABLE ONLY donoengine.category ADD CONSTRAINT category_code_check
    CHECK (code <> '' AND length(code) <= 32);
ALTER TABLE ONLY donoengine.category ADD CONSTRAINT category_name_check
    CHECK (name <> '' AND length(name) <= 32);

-- Category Relation

CREATE TABLE donoengine.category_relation (
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    parent_id       bigint,
    child_id        bigint
);

ALTER TABLE ONLY donoengine.category_relation ADD CONSTRAINT category_relation_pkey
    PRIMARY KEY (parent_id, child_id);

ALTER TABLE ONLY donoengine.category_relation ADD CONSTRAINT category_relation_parent_id_fkey
    FOREIGN KEY (parent_id) REFERENCES donoengine.category(id) ON DELETE CASCADE;
ALTER TABLE ONLY donoengine.category_relation ADD CONSTRAINT category_relation_child_id_fkey
    FOREIGN KEY (child_id) REFERENCES donoengine.category(id) ON DELETE CASCADE;

ALTER TABLE ONLY donoengine.category_relation ADD CONSTRAINT category_relation_parent_id_child_id_check
    CHECK (parent_id <> child_id);

-- Tag Type

CREATE TABLE donoengine.tag_type (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    code            text                        NOT NULL,
    name            text                        NOT NULL
);

ALTER TABLE ONLY donoengine.tag_type ADD CONSTRAINT tag_type_code_key
    UNIQUE (code);

ALTER TABLE ONLY donoengine.tag_type ADD CONSTRAINT tag_type_code_check
    CHECK (code <> '' AND length(code) <= 24);
ALTER TABLE ONLY donoengine.tag_type ADD CONSTRAINT tag_type_name_check
    CHECK (name <> '' AND length(name) <= 24);

-- Tag

CREATE TABLE donoengine.tag (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    type_id         bigint                      NOT NULL,
    code            text                        NOT NULL,
    name            text                        NOT NULL
);

ALTER TABLE ONLY donoengine.tag ADD CONSTRAINT tag_type_id_fkey
    FOREIGN KEY (type_id) REFERENCES donoengine.tag_type(id);

ALTER TABLE ONLY donoengine.tag ADD CONSTRAINT tag_type_id_code_key
    UNIQUE (type_id, code);

ALTER TABLE ONLY donoengine.tag ADD CONSTRAINT tag_code_check
    CHECK (code <> '' AND length(code) <= 32);
ALTER TABLE ONLY donoengine.tag ADD CONSTRAINT tag_name_check
    CHECK (name <> '' AND length(name) <= 32);

-- Comic

CREATE TABLE donoengine.comic (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    code            text                        NOT NULL,
    language_id     bigint,
    published_from  timestamp with time zone,
    published_to    timestamp with time zone,
    total_chapter   int,
    total_volume    int,
    nsfw            int,
    nsfl            int,
    additionals     jsonb
);

ALTER TABLE ONLY donoengine.comic ADD CONSTRAINT comic_language_id_fkey
    FOREIGN KEY (language_id) REFERENCES donoengine.language(id);

ALTER TABLE ONLY donoengine.comic ADD CONSTRAINT comic_code_key
    UNIQUE (code);

ALTER TABLE ONLY donoengine.comic ADD CONSTRAINT comic_code_check
    CHECK (length(code) = 8);
ALTER TABLE ONLY donoengine.comic ADD CONSTRAINT comic_nsfw_check
    CHECK (nsfw <= -1 AND nsfw >= 1);
ALTER TABLE ONLY donoengine.comic ADD CONSTRAINT comic_nsfl_check
    CHECK (nsfl <= -1 AND nsfl >= 1);

-- Comic Title

CREATE TABLE donoengine.comic_title (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    comic_id        bigint                      NOT NULL,
    rid             text                        NOT NULL,
    language_id     bigint                      NOT NULL,
    title           text                        NOT NULL,
    synonym         boolean,
    romanized       boolean
);

ALTER TABLE ONLY donoengine.comic_title ADD CONSTRAINT comic_title_comic_id_fkey
    FOREIGN KEY (comic_id) REFERENCES donoengine.comic(id) ON DELETE CASCADE;
ALTER TABLE ONLY donoengine.comic_title ADD CONSTRAINT comic_title_language_id_fkey
    FOREIGN KEY (language_id) REFERENCES donoengine.language(id);

ALTER TABLE ONLY donoengine.comic_title ADD CONSTRAINT comic_title_comic_id_rid_key
    UNIQUE (comic_id, rid);
ALTER TABLE ONLY donoengine.comic_title ADD CONSTRAINT comic_title_comic_id_title_key
    UNIQUE (comic_id, title);

ALTER TABLE ONLY donoengine.comic_title ADD CONSTRAINT comic_title_rid_check
    CHECK (length(rid) = 4);
ALTER TABLE ONLY donoengine.comic_title ADD CONSTRAINT comic_title_title_check
    CHECK (title <> '' AND length(title) < 256);

-- Comic Cover

CREATE TABLE donoengine.comic_cover (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    comic_id        bigint                      NOT NULL,
    rid             text                        NOT NULL,
    website_id      bigint                      NOT NULL,
    relative_url    text                        NOT NULL,
    priority        int
);

ALTER TABLE ONLY donoengine.comic_cover ADD CONSTRAINT comic_cover_comic_id_fkey
    FOREIGN KEY (comic_id) REFERENCES donoengine.comic(id) ON DELETE CASCADE;
ALTER TABLE ONLY donoengine.comic_cover ADD CONSTRAINT comic_cover_website_id_fkey
    FOREIGN KEY (website_id) REFERENCES donoengine.website(id) ON DELETE CASCADE;

ALTER TABLE ONLY donoengine.comic_cover ADD CONSTRAINT comic_cover_comic_id_rid_key
    UNIQUE (comic_id, rid);
ALTER TABLE ONLY donoengine.comic_cover ADD CONSTRAINT comic_cover_comic_id_website_id_relative_url_key
    UNIQUE (comic_id, website_id, relative_url);

ALTER TABLE ONLY donoengine.comic_cover ADD CONSTRAINT comic_cover_rid_check
    CHECK (length(rid) = 4);
ALTER TABLE ONLY donoengine.comic_cover ADD CONSTRAINT comic_cover_relative_url_check
    CHECK (relative_url <> '' AND length(relative_url) <= 128);

-- Comic Synopsis

CREATE TABLE donoengine.comic_synopsis (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    comic_id        bigint                      NOT NULL,
    rid             text                        NOT NULL,
    language_id     bigint                      NOT NULL,
    synopsis        text                        NOT NULL,
    version         text,
    romanized       boolean
);

ALTER TABLE ONLY donoengine.comic_synopsis ADD CONSTRAINT comic_synopsis_comic_id_fkey
    FOREIGN KEY (comic_id) REFERENCES donoengine.comic(id) ON DELETE CASCADE;
ALTER TABLE ONLY donoengine.comic_synopsis ADD CONSTRAINT comic_synopsis_language_id_fkey
    FOREIGN KEY (language_id) REFERENCES donoengine.language(id);

ALTER TABLE ONLY donoengine.comic_synopsis ADD CONSTRAINT comic_synopsis_comic_id_rid_key
    UNIQUE (comic_id, rid);
ALTER TABLE ONLY donoengine.comic_synopsis ADD CONSTRAINT comic_synopsis_comic_id_synopsis_key
    UNIQUE (comic_id, synopsis);

ALTER TABLE ONLY donoengine.comic_synopsis ADD CONSTRAINT comic_synopsis_rid_check
    CHECK (length(rid) = 4);
ALTER TABLE ONLY donoengine.comic_synopsis ADD CONSTRAINT comic_synopsis_synopsis_check
    CHECK (synopsis <> '' AND length(synopsis) <= 2048);
ALTER TABLE ONLY donoengine.comic_synopsis ADD CONSTRAINT comic_synopsis_version_check
    CHECK (version <> '' AND length(version) <= 12);

-- Comic External

CREATE TABLE donoengine.comic_external (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    comic_id        bigint                      NOT NULL,
    rid             text                        NOT NULL,
    website_id      bigint                      NOT NULL,
    relative_url    text,
    official        boolean
);

ALTER TABLE ONLY donoengine.comic_external ADD CONSTRAINT comic_external_comic_id_fkey
    FOREIGN KEY (comic_id) REFERENCES donoengine.comic(id) ON DELETE CASCADE;
ALTER TABLE ONLY donoengine.comic_external ADD CONSTRAINT comic_external_website_id_fkey
    FOREIGN KEY (website_id) REFERENCES donoengine.website(id) ON DELETE CASCADE;

ALTER TABLE ONLY donoengine.comic_external ADD CONSTRAINT comic_external_comic_id_rid_key
    UNIQUE (comic_id, rid);
ALTER TABLE ONLY donoengine.comic_external ADD CONSTRAINT comic_external_comic_id_website_id_relative_url_key
    UNIQUE (comic_id, website_id, relative_url);

ALTER TABLE ONLY donoengine.comic_external ADD CONSTRAINT comic_external_rid_check
    CHECK (length(rid) = 4);
ALTER TABLE ONLY donoengine.comic_external ADD CONSTRAINT comic_external_relative_url_check
    CHECK (relative_url <> '' AND length(relative_url) <= 128);

-- Comic Category

CREATE TABLE donoengine.comic_category (
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    comic_id        bigint,
    category_id     bigint
);

ALTER TABLE ONLY donoengine.comic_category ADD CONSTRAINT comic_category_pkey
    PRIMARY KEY (comic_id, category_id);

ALTER TABLE ONLY donoengine.comic_category ADD CONSTRAINT comic_category_comic_id_fkey
    FOREIGN KEY (comic_id) REFERENCES donoengine.comic(id) ON DELETE CASCADE;
ALTER TABLE ONLY donoengine.comic_category ADD CONSTRAINT comic_category_category_id_fkey
    FOREIGN KEY (category_id) REFERENCES donoengine.category(id) ON DELETE CASCADE;

-- Comic Tag

CREATE TABLE donoengine.comic_tag (
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    comic_id        bigint,
    tag_id          bigint
);

ALTER TABLE ONLY donoengine.comic_tag ADD CONSTRAINT comic_tag_pkey
    PRIMARY KEY (comic_id, tag_id);

ALTER TABLE ONLY donoengine.comic_tag ADD CONSTRAINT comic_tag_comic_id_fkey
    FOREIGN KEY (comic_id) REFERENCES donoengine.comic(id) ON DELETE CASCADE;
ALTER TABLE ONLY donoengine.comic_tag ADD CONSTRAINT comic_tag_tag_id_fkey
    FOREIGN KEY (tag_id) REFERENCES donoengine.tag(id) ON DELETE CASCADE;

-- Comic Relation Type

CREATE TABLE donoengine.comic_relation_type (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    code            text                        NOT NULL,
    name            text                        NOT NULL
);

ALTER TABLE ONLY donoengine.comic_relation_type ADD CONSTRAINT comic_relation_type_code_key
    UNIQUE (code);

ALTER TABLE ONLY donoengine.comic_relation_type ADD CONSTRAINT comic_relation_type_code_check
    CHECK (code <> '' AND length(code) <= 24);
ALTER TABLE ONLY donoengine.comic_relation_type ADD CONSTRAINT comic_relation_type_name_check
    CHECK (name <> '' AND length(name) <= 24);

-- Comic Relation

CREATE TABLE donoengine.comic_relation (
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    type_id         bigint,
    parent_id       bigint,
    child_id        bigint
);

ALTER TABLE ONLY donoengine.comic_relation ADD CONSTRAINT comic_relation_pkey
    PRIMARY KEY (type_id, parent_id, child_id);

ALTER TABLE ONLY donoengine.comic_relation ADD CONSTRAINT comic_relation_type_id_fkey
    FOREIGN KEY (type_id) REFERENCES donoengine.comic_relation_type(id);
ALTER TABLE ONLY donoengine.comic_relation ADD CONSTRAINT comic_relation_parent_id_fkey
    FOREIGN KEY (parent_id) REFERENCES donoengine.comic(id) ON DELETE CASCADE;
ALTER TABLE ONLY donoengine.comic_relation ADD CONSTRAINT comic_relation_child_id_fkey
    FOREIGN KEY (child_id) REFERENCES donoengine.comic(id) ON DELETE CASCADE;

ALTER TABLE ONLY donoengine.comic_relation ADD CONSTRAINT comic_relation_parent_id_child_id_check
    CHECK (parent_id <> child_id);

-- Comic Chapter

CREATE TABLE donoengine.comic_chapter (
    id              bigint                      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at      timestamp with time zone    NOT NULL DEFAULT timezone('UTC', now()),
    updated_at      timestamp with time zone,

    comic_id        bigint                      NOT NULL,
    chapter         text                        NOT NULL,
    version         text,
    volume          text,
    released_at     timestamp with time zone    NOT NULL
);

ALTER TABLE ONLY donoengine.comic_chapter ADD CONSTRAINT comic_chapter_comic_id_fkey
    FOREIGN KEY (comic_id) REFERENCES donoengine.comic(id) ON DELETE CASCADE;

ALTER TABLE ONLY donoengine.comic_chapter ADD CONSTRAINT comic_chapter_comic_id_chapter_version_key
    UNIQUE NULLS NOT DISTINCT (comic_id, chapter, version);

ALTER TABLE ONLY donoengine.comic_chapter ADD CONSTRAINT comic_chapter_chapter_check
    CHECK (chapter <> '' AND length(chapter) <= 64);
ALTER TABLE ONLY donoengine.comic_chapter ADD CONSTRAINT comic_chapter_version_check
    CHECK (version <> '' AND length(version) <= 32);
ALTER TABLE ONLY donoengine.comic_chapter ADD CONSTRAINT comic_chapter_volume_check
    CHECK (volume <> '' AND length(volume) <= 24);

-- +goose Down

DROP TABLE donoengine.comic_chapter;
DROP TABLE donoengine.comic_relation;
DROP TABLE donoengine.comic_relation_type;
DROP TABLE donoengine.comic_tag;
DROP TABLE donoengine.comic_category;
DROP TABLE donoengine.comic_external;
DROP TABLE donoengine.comic_synopsis;
DROP TABLE donoengine.comic_cover;
DROP TABLE donoengine.comic_title;
DROP TABLE donoengine.comic;
DROP TABLE donoengine.tag;
DROP TABLE donoengine.tag_type;
DROP TABLE donoengine.category_relation;
DROP TABLE donoengine.category;
DROP TABLE donoengine.category_type;
DROP TABLE donoengine.website;
DROP TABLE donoengine.language;
DROP SCHEMA donoengine;
