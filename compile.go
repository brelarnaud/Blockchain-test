package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}
// poussé
type Product struct {
	Ref         string  `json:"ref"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Critical    int     `json:"critical"`
	Provision   int     `json:"provision"`
}

type Order struct {
	Ref        string    `json:"ref"`
	UserHash   string    `json:"user"`
	Products   []Product `json:"products"`
	Quantities []int     `json:"quantities"`
	TotalPrice float64   `json:"totalprice"`
	TrackingID string    `json:"trackingid"`
	State      int       `json:"state"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	var err error

	err = stub.PutState("productsLength", []byte("0"))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("ordersLength", []byte("0"))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("usersLength", []byte("0"))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "addProduct" {
		return t.addProduct(stub, args)
	} else if function == "addOrder" {
		return t.addOrder(stub, args)
	} else if function == "addUser" {
		return t.addUser(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}




/************************************************************************************************************/
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}



/**********************************************************************************************************/
func (t *SimpleChaincode) addProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	var err error
	var product Product

	productsLengthAsBytes, err := stub.GetState("productsLength")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for productsLength\"}"
		return nil, errors.New(jsonResp)
	}

	product.Ref = args[0]
	product.Description = args[1]
	product.Price, err = strconv.ParseFloat(args[2], 64)
	product.Quantity, err = strconv.Atoi(args[3])
	product.Critical, err = strconv.Atoi(args[4])
	product.Provision = 0

	productAsBytes, err := json.Marshal(product)
	productsLength := string(productsLengthAsBytes)
	err = stub.PutState("product"+productsLength, productAsBytes)
	if err != nil {
		return nil, err
	}

	count, err := strconv.Atoi(productsLength)
	count++
	err = stub.PutState("productsLength", []byte(strconv.Itoa(count)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}



/***********************************************************************************************************/
func (t *SimpleChaincode) setProvision(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	var index, value string
	var product Product

	index = args[0]
	value = args[1]

	prouctsAsBytes, err := stub.GetState("product" + index)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for : product" + index + "\"}"
		return nil, errors.New(jsonResp)
	}

	err = json.Unmarshal(prouctsAsBytes, &product)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to unmarshal for : product" + index + "\"}"
		return nil, errors.New(jsonResp)
	}

	product.Provision, err = strconv.Atoi(value)
	prouctsAsBytes, err = json.Marshal(product)
	err = stub.PutState("product"+index, prouctsAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}




/******************************************************************************************************/
func (t *SimpleChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var err error
	var userLogin, userPassword, userHash string

	userLogin = args[0]
	userPassword = args[1]
	userHash = args[2]

	usersLengthAsBytes, err := stub.GetState("usersLength")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for usersLength\"}"
		return nil, errors.New(jsonResp)
	}

	err = stub.PutState(userLogin+"@"+userPassword, []byte(string(userHash)))
	if err != nil {
		return nil, err
	}

	usersLength := string(usersLengthAsBytes)
	count, err := strconv.Atoi(usersLength)
	count++
	err = stub.PutState("usersLength", []byte(strconv.Itoa(count)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}







/***********************************************************************************************************/
func (t *SimpleChaincode) addOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}
	fmt.Println("args[0] : " + args[0])
	fmt.Println("args[1] : " + args[1])
	fmt.Println("args[2] : " + args[2])
	fmt.Println("args[3] : " + args[3])

	var err error
	var order Order

	userHashAsBytes, err := stub.GetState(args[0])
	fmt.Println("userHashAsBytes:")
	fmt.Println(userHashAsBytes)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	ordersLengthAsBytes, err := stub.GetState("ordersLength")
	fmt.Println("ordersLenghtAsBytes:")
	fmt.Println(ordersLengthAsBytes)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for ordersLength\"}"
		return nil, errors.New(jsonResp)
	}
	ordersLength := string(ordersLengthAsBytes)
	fmt.Println("ordersLength:")
	fmt.Println(ordersLength)
	count, err := strconv.Atoi(ordersLength)
	fmt.Println("currenCount:")
	fmt.Println(count)
	count++
	fmt.Println("incrementCount:")
	fmt.Println(count)

	order.Ref = strconv.Itoa(count)
	order.UserHash = string(userHashAsBytes)
	err = json.Unmarshal([]byte(args[1]), &order.Products)
	fmt.Println("order.Products:")
	fmt.Println(order.Products)
	fmt.Println("err unmarshal args[1]:")
	fmt.Println(err)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to unmarshal for :\n" + args[1] + "\"}"
		return nil, errors.New(jsonResp)
	}
	err = json.Unmarshal([]byte(args[2]), &order.Quantities)
	fmt.Println("order.Quantities:")
	fmt.Println(order.Quantities)
	fmt.Println("err unmarshal args[2]:")
	fmt.Println(err)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to unmarshal for :\n" + args[3] + "\"}"
		return nil, errors.New(jsonResp)
	}
	order.TotalPrice, err = strconv.ParseFloat(args[3], 64)
	order.TrackingID = ""
	order.State = 1

	ordersAsBytes, err := json.Marshal(order)
	fmt.Println("ordersAsBytes:")
	fmt.Println(ordersAsBytes)
	fmt.Println("err marshal ordersAsBytes:")
	fmt.Println(err)
	err = stub.PutState("order"+ordersLength, ordersAsBytes)
	fmt.Println("key:")
	fmt.Println("order" + ordersLength)
	fmt.Println("err putting stat:")
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	err = stub.PutState("ordersLength", []byte(strconv.Itoa(count)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}


/***************************************************************************************************/

func (t *SimpleChaincode) MajProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	var i,j int
	var err error
	var product Product
	var order Order
	var qt int 
	var reforder string
	var refproduct string
	var critiproduct int
	var qtproduct int
	var qtfinal int


	fmt.Println("args[0]:")
	fmt.Println(order.Products)

	orderAsbytes := []byte(args[0])
	err = json.Unmarshal(orderAsbytes, &order)
	fmt.Println("ordersAsBytes:")
	fmt.Println(orderAsbytes)
	fmt.Println("err marshal err:")
	fmt.Println(err)

	productsLengthAsBytes, err := stub.GetState("productsLength")
	productsLength, err := strconv.Atoi(string(productsLengthAsBytes)) 

	//On recupère la longueur de la liste
	ordersLengthAsBytes, err := stub.GetState("ordersLength")
	ordersLength, err := strconv.Atoi(string(ordersLengthAsBytes))

	fmt.Println("ordersLength:")
	fmt.Println(ordersLength)
	

	//On prend chaque élément de la liste un par un
		for i = 0; i < len(order.Products); i++ {
		productsAsBytes, err := stub.GetState("product"+ strconv.Itoa(i))
		err = json.Unmarshal(productsAsBytes, &product)
		_ = err

		refproduct = product.Ref
		qt = product.Quantity
		fmt.Println("refproduct:")
		fmt.Println(refproduct)
		fmt.Println("qt:")
		fmt.Println(qt)

			for j = 0; j < productsLength; j++  {
			productAsBytes, err := stub.GetState("product"+ strconv.Itoa(j))
			err = json.Unmarshal(productsAsBytes, &product)
			_ = productAsBytes
			_ = err
				if reforder == product.Ref {

				//on  recupère la quantité et seuil du produit égal à reforder 
				critiproduct = product.Critical
				qtproduct = product.Quantity
				fmt.Println("critiproduct:")
				fmt.Println(critiproduct)
				fmt.Println("qtproduct:")
				fmt.Println(qtproduct)

					if qtproduct > qt {
					//on vérifie si la quantité en stock supérieur a la quantité demandée
					qtfinal = qtproduct - qt
					fmt.Println("qtfinal:")
					fmt.Println(qtfinal)
					//ici envoyé qtfinal a la place de dans product.quantity
					product.Quantity = qtfinal
		
						//sinon envoie erreur

						//Condition réapprovisionnement
						if critiproduct < qtfinal {
						fmt.Println("Event : commande en cours sur le produit X")
						//déclancher la commande


						}
					} 
		
		

				}
			}
		}
	return nil, nil

}
