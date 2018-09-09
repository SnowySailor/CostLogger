package main

import (
    "testing"
)

func TestValidatePassword(t *testing.T) {
    tables := []struct {
        x string
        y string
        z bool
    } {
        {"123", "$2b$10$V0v2iM96ALGRq4AnqZ0pz.gMz8h9Nw9S3kDQpeaXA7pZlHWkKJhrS", true},
        {"testing123", "$2y$12$YUIRFclPbegtvFxPHrYbqu5htrrtlH1f47xRo3Gpjh//Go02SWnfS", true},
        {"Testing123", "$2y$12$VlVtK3DVg37vFcgENttv1ugr4L.0NeihdN0DeHtBllZ9fyn7DdXCK", true},
        {"¶¶••ªªºº¡¡¡¡™™", "$2y$10$oTlwZgBTkcjwGebMwTPMLO.sSRdsCAjEcomT9IGo1eNBLidAuoABa", true},
        {"¶¶••ªªºº¡¡¡¡™™", "$2y$11$oTlwZgBTkcjwGebMwTPMLO.sSRdsCAjEcomT9IGo1eNBLidAuoABa", false},
        {"¶¶••ªªºº¡¡¡¡™™", "$2y$11$oTlwZgBTkcjwGebMwTPMLO.sSRdsCAjEcomT9IGo1eNBLidAuoABc", false},
        {"漢字", "$2y$12$m0MDRc5ReNdxiyCcidSp5.258qXWoqx7Svs2c4Icjfkp8XitSSewy", true},
        {"汉字", "$2y$10$ozYTs.qWGm.2250w28CYVuWj8.iA1aCQGNyfeJMcsBInPvUVYC8CC", true},
        {"汉字", "$2y$10$ozYTs.qWGm.2250w28CYVuWj8.iA1aCQGNyfeJMcsBInPvUVYC8CD", false},
        {"", "$2y$14$cWA4rxVKQYe//8vYHL23NOWjo/ivoy1usL93chEI.NpdPvV3Smk3u", true},
        {"whataverytrulylongpasswordthishereisandIamgoingtogogetsomeicecreamlaterbecauseIamvery,veryhungry!Iloveice cream and cake", "$2y$10$SCLvig6RnNeEfJ/tgPoaFO3p9QCOoSKEQsJ5kzmNK3jt4BPiY3W5K", true},
    }

    for _, table := range tables {
        isValid := validatePassword(table.x, table.y)
        if isValid != table.z {
            if table.z {
                t.Errorf("Got False, expected True for password '%s' and hash '%s'", table.x, table.y)
            } else {
                t.Errorf("Got True, expected False for password '%s' and hash '%s'", table.x, table.y)
            }
        }
    }
}

func TestHashPassword(t *testing.T) {
    passwords := []string{
        "testing123", "123", "hello", "•••¡¡™™∆˜®˜®ü", "¶¶••ªªºº¡¡¡¡™™", "漢字", "",
        "whataverytrulylongpasswordthishereisandIamgoingtogogetsomeicecreamlaterbecauseIamvery,veryhungry!Iloveice cream and cake",
        "NULL",
    }

    for _, password := range passwords {
        hash    := hashPassword(password)
        isValid := validatePassword(password, hash)
        if !isValid {
            t.Errorf("Hashed password '%s' and got '%s', but validatePassword returned False.", password, hash)
        }
    }
}

func TestRemoveLeadingSlash(t *testing.T) {
    strings := []struct {
        x string
        y string
    } {
        {"", ""},
        {"/", ""},
        {"/ ", " "},
        {"/yes", "yes"},
        {"yes/no", "yes/no"},
        {"/yes/no/maybe123", "yes/no/maybe123"},
        {"/begin/and/end/", "begin/and/end/"},
        {"/begin?yes=1&no=0", "begin?yes=1&no=0"},
    }

    for _, str := range strings {
        if removeLeadingSlash(str.x) != str.y {
            t.Errorf("Expected '%s' got '%s'", str.y, removeLeadingSlash(str.x))
        }
    }
}
