package cyk

import (
	"errors"
	"fmt"
	"strings"
)

//Grammar is present the CFG grammars
//ex:  S-> AB
//     LeftSymbol: S
//     RightSymbol: "AB"
type Grammar struct {
	LeftSymbol  string
	RightSymbol string
}

//We use matrix to store our CYK result and maniplate it
//ex: X_11 -> A ==> map[MatrixIndicator{ X_axi: 1, Y_axi: 1}] = "A"
type MatrixIndicator struct {
	X_axi int
	Y_axi int
}

//Using map to handle matrix result
type MatrixResult map[MatrixIndicator]string

type CYK struct {
	Grammars    []Grammar
	CYKResult   MatrixResult
	InputString string
	StartSymbol string
}

func NewCYK(startSymbol string) *CYK {
	newCYK := &CYK{StartSymbol: startSymbol}
	newCYK.CYKResult = make(map[MatrixIndicator]string)
	return newCYK
}

// Find terminal assign variable
// ex: A->a  using `a` find A
func (c *CYK) findTerminalAssign(terminal string) string {
	var retList string
	for _, targetG := range c.Grammars {
		if strings.Contains(targetG.RightSymbol, terminal) {
			retList = fmt.Sprintf("%s%s", retList, targetG.LeftSymbol)
		}
	}

	return retList
}

// Find variable assign, difference wit find terminal it need equal not contains
// S -> AB, only use "AB" -> "S"
func (c *CYK) findVariableAssign(symbol string) string {
	var retSlice string
	for _, targetG := range c.Grammars {
		//fmt.Println(" grammarR=", targetG.RightSymbol, " symbol=", symbol)
		if symbol == targetG.RightSymbol {
			retSlice = fmt.Sprintf("%s%s", retSlice, targetG.LeftSymbol)
			//fmt.Println("get Left=", retSlice)
		}
	}
	return retSlice
}

//To eval if string is terminal not variable
func (c *CYK) isTerminal(testString string) bool {
	return testString == strings.ToLower(testString)
}

//Insert grammar in this CYK
// Ex: S->{AB}  InputGrammar("S", "AB")
// Please note: Uppercase is variable, Lowercase is terminal
func (c *CYK) InputGrammar(leftSymbol string, rightSymbols string) {
	newGrammar := Grammar{LeftSymbol: leftSymbol, RightSymbol: rightSymbols}
	c.Grammars = append(c.Grammars, newGrammar)
}

// Eval takes input to run CYK algorithm and eval if this input can be generated by this CFG
func (c *CYK) Eval(input string) bool {
	c.runCYK(input)
	return c.evalCYKResult()
}

func (c *CYK) getResultMatrix(x, y int) (string, error) {
	val, ok := c.CYKResult[MatrixIndicator{X_axi: x, Y_axi: y}]
	if ok {
		return val, nil
	} else {
		//fmt.Println("index x=", x, " y=", y, " is not exist!.")
		return "", errors.New("Not exist!")
	}
}

func (c *CYK) setResultMatrix(x, y int, val string) {
	c.CYKResult[MatrixIndicator{X_axi: x, Y_axi: y}] = val
}

// Run CYK algorithm
func (c *CYK) runCYK(input string) {
	c.InputString = input
	//Start to calculate X_11, X_22, X_33
	for i := 0; i < len(input); i++ {
		variable := c.findTerminalAssign(string(input[i]))
		c.setResultMatrix(i, i, variable)
	}

	//start triangle calculate
	for loop := 1; loop <= len(c.InputString); loop++ {
		for i := 0; i < len(c.InputString)-loop; i++ {
			j := i + loop
			//fmt.Println("i=", i, " j=", j)
			var totalTargets []string
			for k := 1; k <= j; k++ {

				firstVal, _ := c.getResultMatrix(i, i+k-1)
				secondVal, _ := c.getResultMatrix(i+k, j)
				products := arrayProduction(firstVal, secondVal)
				for _, v := range products {
					totalTargets = append(totalTargets, v)
				}
				//fmt.Println("total =", totalTargets)
			}

			var result string
			for _, symbol := range totalTargets {
				//fmt.Println("i=", i, " j=", j, " symbol=", symbol)
				targetSymbol := c.findVariableAssign(symbol)
				if !strings.Contains(result, targetSymbol) {
					result = fmt.Sprintf("%s%s", result, targetSymbol)
				}
			}

			c.setResultMatrix(i, j, result)
		}
	}
}

func arrayProduction(str1, str2 string) []string {
	var ret []string
	for i := 0; i < len(str1); i++ {
		for j := 0; j < len(str2); j++ {
			restr := fmt.Sprintf("%s%s", string(str1[i]), string(str2[j]))
			ret = append(ret, restr)
		}
	}
	return ret
}

// Eval CYK result and make sure latest CYK Result only contain variable not assign to terminal
// ex: latest result is "S" which S->AB
func (c *CYK) evalCYKResult() bool {
	finalResult, err := c.getResultMatrix(0, len(c.InputString)-1)
	if err != nil {
		return false
	}

	//fmt.Println("final:", finalResult)
	if strings.Contains(finalResult, c.StartSymbol) {
		return true
	}

	return false
}

// Print out the triangle result on CYK
func (c *CYK) PrintResult() {
	if len(c.CYKResult) == 0 {
		fmt.Println("We still not calculate CYK or no result...")
		return
	}

	fmt.Printf("1:")
	for i := 0; i < len(c.InputString); i++ {
		c.printResultMatrixElement(i, i)
	}
	fmt.Printf("\n")

	lineIndex := 2
	for loop := 1; loop < len(c.InputString); loop++ {

		fmt.Printf("%d:", lineIndex)
		for i := 0; i < len(c.InputString)-loop; i++ {
			j := i + loop
			c.printResultMatrixElement(i, j)
		}
		fmt.Printf("\n")
		lineIndex = lineIndex + 1
	}
}

func (c *CYK) printResultMatrixElement(i, j int) {
	fmt.Printf("\tX%d%d:{", i+1, j+1)

	results, err := c.getResultMatrix(i, j)
	if err != nil {
		fmt.Println("Empty result")
		return
	}
	for index, _ := range results {
		fmt.Printf("%s", string(results[index]))
		if index < len(results)-1 {
			fmt.Printf(",")
		}
	}
	fmt.Printf("}")
}
