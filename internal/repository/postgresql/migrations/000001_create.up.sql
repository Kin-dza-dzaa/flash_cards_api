CREATE TABLE IF NOT EXISTS word_translation(
    word                TEXT                                                                NOT NULL CHECK(word != ''),
    trans_data          JSONB                                                               NOT NULL,
    PRIMARY KEY (word)
);

CREATE TABLE IF NOT EXISTS user_collection(
    user_id                                     TEXT                                        NOT NULL,
    word                                        TEXT                                        NOT NULL CHECK(word != ''),
    collection_name                             TEXT                                        NOT NULL,
    time_diff                                   INTERVAL                                    NOT NULL,
    last_repeat                                 TIMESTAMP                                   NOT NULL,
    FOREIGN KEY (word) REFERENCES word_translation(word),
    UNIQUE(user_id, word, collection_name)
);
