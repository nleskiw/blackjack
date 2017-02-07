/*
Blackjack

One Player Simple Blackjack
	* Hit and Stand only, (no double, split or insurance)
	* Dealer stands on soft 17

* Are there less than 17 cards in the deck?
	Note:
		AFAIK, the most cards for 2 players possible is 16:
		player hits to 21 ( A A A A 2 2 2 2 3 3 3 )
		dealer hits to 19 ( 3 4 4 4 4 )
	- True:
		* Shuffle deck.

* Is player wallet <= 1?
	- True:
		* Exit Game, player loses
* Bet = 0
* While player bet is 0
	* Ask player how much to bet
	* convert to int
	* Error converting to int?
		- True:
			* bet = 0
	* Is bet > player wallet?
		- True:
			* bet = 0
	* Is bet less than 0?
		- True:
			* Bet = 0
* Take Bet from player wallet
* Deal 2 cards to Dealer and player
* Does dealer have Ace + 10 K Q J ?
		* Does player have Ace + 10 K Q J?
			- True:
				* Show all cards
				* Push.
				* Return bet.
				* Next hand
			- False:
				* Show all cards.
				* Loss.
				* Next hand
* Does player have Ace + 10 K Q J?
	- True:
		* Player win.
		* Bet = Bet * 1.5, round down to lowest integer
		* Put winnings in player wallet.
		* Next Hand.
* Show 1 dealer card to player
* Show both player cards
* While player total < 21
	* Ask player "Hit or Stand"
	* If hit:
		* Deal card to player
		* Next loop
	* If stand:
		* exit loop
* Is player total > 21?
	- True:
		* Loss.
		* Next hand.
* Show dealer's 2nd card
* While dealer total < 17
	* Deal card to dealer
* Is dealer total > 21
	- True:
		* Win
		* Double bet and return to player's wallet
		* next hand
* Is player total == dealer total?
	- True:
		* Push
		* Return bet
		* Next hand.
* Is player total > dealer total?
	- True:
		* Win
		* Double bet and return to player wallet.
		* Return bet
		* Next Hand
	- False:
		* Loss
		* Next hand.

*phew*
That's a lot of rules.

*/

package main

import (
	"fmt"
	//"os"
	"strconv"

	"github.com/nleskiw/goplaycards/deck"
)

// getString gets an arbitrary string from the user with a prompt.
func getString(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

// getInteger gets an arbitrary int from the user with a prompt.
// Retries until a valid integer is entered.
func getInteger(prompt string) int {
	valid := true
	input := getString(prompt)
	integer, err := strconv.Atoi(input)
	if err != nil {
		valid = false
	}
	for valid == false {
		fmt.Println("Can't convert your answer into an integer.")
		input = getString(prompt)
		integer, err = strconv.Atoi(input)
		if err == nil {
			valid = true
		}
	}
	return integer
}

// getBet gets an amount from user and removes the bet from the wallet
// Retries until a valid bet is entered.
func getBet(wallet *float64) int {
	bet := 0
	valid := false
	for valid == false {
		valid = true
		bet = getInteger("How much would you like to bet ($5 increments)? ")
		if bet < 5 {
			valid = false
		}
		if float64(bet) > *wallet {
			valid = false
		}
		if bet%5 != 0 {
			valid = false
		}
		if valid == false {
			fmt.Println("Invalid bet.")
		}
	}
	*wallet = *wallet - float64(bet)
	return bet
}

// getPlayerAction determines what the player will do
// TODO: Implement Double and Split
func getPlayerAction() string {
	validInput := false
	input := ""
	for validInput == false {
		input = getString("[H]it or [S]tand? ")
		if input == "hit" || input == "Hit" || input == "H" || input == "h" {
			validInput = true
			input = "H"
		}
		if input == "stand" || input == "Stand" || input == "S" || input == "s" {
			validInput = true
			input = "S"
		}
		if validInput == false {
			fmt.Println("Invalid option. H to Hit or S to Stand. ")
		}
	}
	return input
}

// handTotal returns the numerical value of a Blackjack hand
func handTotal(hand []deck.Card) int {
	total := 0
	numberOfAces := 0
	for _, card := range hand {
		if card.Value.Name == "Ace" {
			numberOfAces = numberOfAces + 1
		} else {
			if card.Facecard() {
				total = total + 10
			} else {
				total = total + card.Value.Value
			}
		}
	}

	// If there's at least one Ace, deal with it.
	// In multi-shoe decks, there could be many Aces (more than 4) in a hand.
	if numberOfAces > 0 {
		// All but the last Ace must be a one, because 11 + 11 = 22 (bust)
		// This loop shouldn't run if there's only one Ace
		for numberOfAces > 1 {
			total = total + 1
			numberOfAces = numberOfAces - 1
		}
		// There should now only be one Ace
		// if the last Ace being 11 doesn't cause a bust, make it an 11
		if total+11 > 21 {
			total = total + 1
		} else {
			// If 11 doesn't cause a bust, make it worth 11
			total = total + 11
		}
	}
	return total
}

// Returns true if a hand is bust / over 21
func isBust(hand []deck.Card) bool {
	if handTotal(hand) > 21 {
		return true
	}
	return false
}

// Returns true if a hand is a Blacjack (Ace + [10 | K | Q | J])
func isBlackjack(hand []deck.Card) bool {
	// A Blackjack is exactly one Ace and Exactly one 10, K, Q, or A
	if len(hand) != 2 {
		return false
	}
	// In the goplaycards library, the values enumerate from 2 to 14.
	// J = 11, Q = 12, K = 13, Ace = 14
	if hand[0].Value.Name == "Ace" {
		if hand[1].Value.Value >= 10 && hand[1].Value.Value <= 13 {
			return true
		}
	}
	if hand[1].Value.Name == "Ace" {
		if hand[0].Value.Value >= 10 && hand[0].Value.Value <= 13 {
			return true
		}
	}
	return false
}

func printHand(hand []deck.Card) {
	for _, card := range hand {
		fmt.Printf("%s  ", card.ToStr())
	}
	fmt.Printf(" Total: %d\n", handTotal(hand))
}

func printPlayerHand(hand []deck.Card) {
	fmt.Printf("Player Hand: ")
	printHand(hand)
}

func printDealerHand(hand []deck.Card, hideFirst bool) {
	fmt.Printf("Dealer Hand: ")
	if hideFirst {
		fmt.Printf("XX  %s  \n", hand[1].ToStr())
	} else {
		printHand(hand)
	}
}

func printLicense() {
	fmt.Println("goplaycards Copyright (C) 2017  Nicholas Leskiw")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY; for details please visit")
	fmt.Println("<https://www.gnu.org/licenses/gpl-3.0.txt>")
	fmt.Println("This is free software, and you are welcome to redistribute it")
	fmt.Printf("under certain conditions. Please see the GPLv3 license at the URL listed above.\n\n")
}

func main() {
	printLicense()
	wallet := 100.00
	bet := 0
	var d deck.Deck
	d.Initialize()
	d.Shuffle()

	// While the player still has enough money for one bet
	for wallet > 5.0 {
		fmt.Println("=============================================================")
		// Minimum number of cards to play a hand from a single deck in worst case
		fmt.Printf("%d cards left in the deck.\n", d.CardsLeft())
		if d.CardsLeft() < 17 {
			// If there's less than 17, you must shuffle the deck.
			fmt.Println("Shuflling deck...")
			d.Initialize()
			d.Shuffle()
			fmt.Printf("%d Cards left in the deck.\n", d.CardsLeft())
		}
		fmt.Printf("You have %.2f in your wallet.\n", wallet)
		bet = getBet(&wallet)
		fmt.Printf("\nYou bet: %d\n", bet)

		// draw the initial hands
		playerHand, err := d.Draw(2)
		if err != nil {
			panic(err)
		}
		dealerHand, err := d.Draw(2)
		if err != nil {
			panic(err)
		}

		// Handle all Blackjack
		if isBlackjack(dealerHand) || isBlackjack(playerHand) {
			printPlayerHand(playerHand)
			printDealerHand(dealerHand, false)

			if isBlackjack(dealerHand) && isBlackjack(playerHand) {
				fmt.Println("Both the Player and the Dealer have Blackjack. Hand is a push.")
				wallet = wallet + float64(bet)
				continue
			}

			if isBlackjack(dealerHand) {
				fmt.Println("Dealer has Blackjack. Player loses this hand.")
				continue
			}

			if isBlackjack(playerHand) {
				winnings := float64(bet) * 2.5
				fmt.Printf("Player has Blackjack. Player wins $%.2f.\n", winnings)
				wallet = wallet + winnings
				continue
			}
		}

		// Let the player play
		playerDone := false
		for playerDone == false {
			printPlayerHand(playerHand)
			printDealerHand(dealerHand, true)
			action := getPlayerAction()
			if action == "H" {
				drawnCards, err := d.Draw(1)
				if err != nil {
					panic(err)
				}
				playerHand = append(playerHand, drawnCards[0])
				if isBust(playerHand) {
					printPlayerHand(playerHand)
					fmt.Println("Bust")
					playerDone = true
				}
			}
			if action == "S" {
				playerDone = true
			}
		}
		if isBust(playerHand) == true {
			continue
		}
		dealerDone := false
		for dealerDone == false {
			if handTotal(dealerHand) >= 17 {
				dealerDone = true
				continue
			}
			drawnCards, err := d.Draw(1)
			if err != nil {
				panic(err)
			}
			dealerHand = append(dealerHand, drawnCards[0])
		}
		printDealerHand(dealerHand, false)

		if handTotal(dealerHand) > 21 {
			fmt.Println("Dealer busts.  Player wins.")
			wallet = wallet + float64(bet*2)
			continue
		}

		if handTotal(playerHand) > handTotal(dealerHand) {
			fmt.Println("Player's hand beats the dealer. Player wins.")
			wallet = wallet + float64(bet*2)
			continue
		}

		if handTotal(playerHand) == handTotal(dealerHand) {
			fmt.Println("Push.")
			wallet = wallet + float64(bet)
			continue
		}

		if handTotal(playerHand) < handTotal(dealerHand) {
			fmt.Println("Dealer wins. Player loses.")
			continue
		}
		fmt.Println("There was a problem determining the winner.")
	}
	fmt.Println("You're out of money.")
	return
}
