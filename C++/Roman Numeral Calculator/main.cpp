#include <iostream>
#include "romanNumeral.h"

int main() {

    romanNumeral  r1("XLIX");
    romanNumeral  r2("XXXVIII");
    romanNumeral  r3;
    romanNumeral  r4;
    romanNumeral  result;

    int  base10Num;

    cout << "Enter a roman numeral: ";
    cin >> r3;
    cout << "\n";
    // r3 = romanNumeral("XIX");

    cout << "Enter a roman numeral: ";
    cin >> r4;
    cout << "\n";
    // r4 = romanNumeral("III");

    cout << "Enter a base 10 number: ";
    cin >> base10Num;
    cout << "\n";
    // base10Num = 13;

    cout << "r1 = " << r1 << "\n";
    cout << "r2 = " << r2 << "\n";
    cout << "r3 = " << r3 << "\n";
    cout << "r4 = " << r4 << "\n";
    cout << "base10Num = " << base10Num << "\n";

    cout << "\n";

    cout << "r1 + r3 = " << r1 + r3 << "\n";
    cout << "r2 + base10Num = " << r2 + base10Num << "\n";
    cout << "r1 * r4 = " << r1 * r4 << "\n";
    cout << "r3 * 6 = " << r3 * 6 << "\n";
    cout << "r1 - r2 = " << r1 - r2 << "\n";
    cout << "r1 - base10Num = " << r1 - base10Num << "\n";
    cout << "100 - r1 = " << 100 - r1 << "\n";
    cout << "2 * r2 = " << 2 * r2 << "\n";
    cout << "base10Num + r2 = " << base10Num + r2 << "\n";

}