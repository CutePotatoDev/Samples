#ifndef NUMERAL_H
#define NUMERAL_H

using namespace std;
#include <string>

class  romanNumeral {

    public:
        romanNumeral ();
        romanNumeral(string);
        romanNumeral  operator +( const  romanNumeral &)  const;
        romanNumeral  operator +(int) const;
        romanNumeral  operator -( const  romanNumeral &)  const;
        romanNumeral  operator -(int) const;
        romanNumeral  operator *( const  romanNumeral &)  const;
        romanNumeral  operator *(int) const;
        friend  romanNumeral  operator +(int , const  romanNumeral &);
        friend  romanNumeral  operator -(int , const  romanNumeral &);
        friend  romanNumeral  operator *(int , const  romanNumeral &);
        friend  ostream& operator <<(ostream&, const  romanNumeral &);
        friend  istream& operator >>(istream&, romanNumeral &);
    
    private:
        string  rNum;
};

#endif