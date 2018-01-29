# Address Validator

Taking addresses via a simple form with free format fields address 1, 2 and 3 will result in address that come through with bad formatting, this package looks up the post code and returns all the addresses associated with it, for the lookup address and each address returned via the postcode lookup will be tokenised and scored.

# Scoring

The scoring is two stage first work out how many fields match, then work out the how many missing fields for each comparison, the less missing fields the better.

# Address data

The address data is based on a response from [Ideal Postcodes](https://ideal-postcodes.co.uk/)
