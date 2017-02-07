Blackjack
=========

Golang implementation of Blackjack (21)

This was written as a Golang exercise to get more familiar with the language.

Table rules:
------------

* Basic play only (Hit / Stand) 
* Dealer stands on soft 17.
* Single deck
* Reshuffle if less than 17 cards in deck.
* Bet in $5 increments only
* Player starts with $100

AFAIK 17 is the most cards you'd need with one deck:
Player (A A A A 2 2 2 2 3 3 3) Dealer (3 4 4 4 4)

Requires https://github.com/nleskiw/goplaycards

Report bugs and/or style suggestions via Github issues.
