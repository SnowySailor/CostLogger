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
        {""    , "$2y$14$cWA4rxVKQYe//8vYHL23NOWjo/ivoy1usL93chEI.NpdPvV3Smk3u", true},
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

    for _, test := range strings {
        if removeLeadingSlash(test.x) != test.y {
            t.Errorf("Expected '%s' got '%s'", test.y, removeLeadingSlash(test.x))
        }
    }
}

func TestPadString(t *testing.T) {
    values := []struct {
        x string
        y int
        z string
        zz bool
        xx string
    } {
        {""     , 6 , "abc"  , false, "abcabc"},
        {"0"    , 4 , "0"    , false, "0000"},
        {"01"   , 10, "a"    , true , "aaaaaaaa01"},
        {"hello", 10, "there", false, "hellothere"},
        {"hello", 11, "there", false, "hellothere"},
        {"hello", 15, "there", false, "hellotherethere"},
        {""     , 0 , ""     , false, ""},
    }
    for _, test := range values {
        if padString(test.x, test.y, test.z, test.zz) != test.xx {
            t.Errorf("Expected '%s' got '%s'", test.xx, padString(test.x, test.y, test.z, test.zz))
        }
    }
}

func TestFlintToString(t *testing.T) {
    values := []struct {
        x  flint
        y  int
        yy int
        z  string
    } {
        {456  , 2 , 2 , "4.56"},
        {0    , 2 , 2 , "0.00"},
        {01   , 2 , 2 , "0.01"},
        {10001, 2 , 2 , "100.01"},
        {9933 , 2 , 2 , "99.33"},
        {100  , 0 , 2 , "100.00"},
        {1001 , 0 , 2 , "1001.00"},
        {1001 , 4 , 4 , "0.1001"},
        {1    , 10, 10, "0.0000000001"},
        {1000 , 4 , 4 , "0.1000"},
        {1000 , 2 , 0 , "10"},
        {1001 , 2 , 1 , "10.0"},
    }
    for _, test := range values {
        if test.x.FlintToString(test.y, test.yy) != test.z {
            t.Errorf("Expected '%s' got '%s'", test.z, test.x.FlintToString(test.y, test.yy))
        }
    }
}

func TestTrim(t *testing.T) {
    values := []struct {
        x string
        y int
        z string
    } {
        {""   ,  0, ""},
        {""   , -1, ""},
        {"a"  ,  0, ""},
        {"abc",  5, "abc"},
        {"abc",  3, "abc"},
        {"abcde", 0, ""},
        {"abcde", 1, "a"},
    }
    for _, test := range values {
        if trim(test.x, test.y) != test.z {
            t.Errorf("Expected '%s' got '%s'", test.z, trim(test.x, test.y))
        }
    }
}