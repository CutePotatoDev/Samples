#include "romanNumeral.h"

// New name space for functions.
namespace Roman {

    // Convert one char.
    int toInt(char c) {
        switch (c) {
            case 'I':  return 1;
            case 'V':  return 5;
            case 'X':  return 10;
            case 'L':  return 50;
            case 'C':  return 100;
            case 'D':  return 500;
            case 'M':  return 1000;
        }
        // Any other character don't have impact.
        return 0;
    }

    // Convert string.
    int toInt(const string& s) {
        int retval = 0, pvs = 0;

        for (auto pc = s.rbegin(); pc != s.rend(); ++pc) {
            const int inc = Roman::toInt(*pc);
            retval += inc < pvs ? -inc : inc;
            pvs = inc;
        }

        return retval;
    }


    struct romandata_t { 
        int value; 
        char const* numeral; 
    };

    static romandata_t const romandata[] = {   
        1000, "M",
        900, "CM",
        500, "D",
        400, "CD",
        100, "C",
        90, "XC",
        50, "L",
        40, "XL",
        10, "X",
        9, "IX",
        5, "V",
        4, "IV",
        1, "I",
        0, NULL 
    };

    // Convert int to roman.
    string toRoman(int value) {
        string result;
      
        for (romandata_t const* current = romandata; current->value > 0; ++current) {
            while (value >= current->value) {
                result += current->numeral;
                value  -= current->value;
            }
        }
        
        return result;
    }

}


// Constructors.
romanNumeral::romanNumeral() {
    rNum = "";
}

romanNumeral::romanNumeral(string str) {
    rNum = str;
}

// IO operators.
istream& operator >> (istream& is, romanNumeral& numeral) {
    is >> numeral.rNum;
    return is;
}

ostream& operator <<(ostream& os, const  romanNumeral& numeral) {
    os << numeral.rNum;
    return os;
}


romanNumeral romanNumeral::operator +(const romanNumeral& numeral) const {
    return romanNumeral(Roman::toRoman(Roman::toInt(this->rNum) + Roman::toInt(numeral.rNum)));
}

romanNumeral romanNumeral::operator +(int num) const {
    return romanNumeral(Roman::toRoman(Roman::toInt(this->rNum) + num));
}

romanNumeral romanNumeral::operator *(const romanNumeral& numeral) const {
    return romanNumeral(Roman::toRoman(Roman::toInt(this->rNum) * Roman::toInt(numeral.rNum)));
}

romanNumeral romanNumeral::operator *(int num) const {
    return romanNumeral(Roman::toRoman(Roman::toInt(this->rNum) * num));
}

romanNumeral romanNumeral::operator -(const romanNumeral& numeral) const {
    return romanNumeral(Roman::toRoman(Roman::toInt(this->rNum) - Roman::toInt(numeral.rNum)));
}

romanNumeral romanNumeral::operator -(int num) const {
    return romanNumeral(Roman::toRoman(Roman::toInt(this->rNum) - num));
}

romanNumeral  operator -(int num, const  romanNumeral& numeral) {
    return romanNumeral(Roman::toRoman(num - Roman::toInt(numeral.rNum)));
}

romanNumeral  operator *(int num, const  romanNumeral& numeral) {
    return romanNumeral(Roman::toRoman(num * Roman::toInt(numeral.rNum)));
}

romanNumeral  operator +(int num, const  romanNumeral& numeral) {
    return romanNumeral(Roman::toRoman(num + Roman::toInt(numeral.rNum)));
}