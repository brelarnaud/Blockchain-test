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

type Collis struct {
	Dimension string  `json:"dimension"`
	Poids     float64 `json:"poids"`
}

type Product struct {
	Ref         string  `json:"ref"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Critical    int     `json:"critical"`
	Provision   int     `json:"provision"`
}

type Order struct {
	Ref         string    `json:"ref"`
	ClientHash  string    `json:"clienthash"`
	CarrierHash string    `json:"carrierhash"`
	Products    []Product `json:"products"`
	Quantities  []int     `json:"quantities"`
	TotalPrice  float64   `json:"totalprice"`
	Collis      Collis    `json:"collis"`
	TrackingID  string    `json:"trackingid"`
	State       int       `json:"state"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

/***************************Init resets all the things********************************/

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	var err error

	err = stub.PutState("productsLength", []byte("0"))
	fmt.Print("PutState: ")
	fmt.Println(err)

	err = stub.PutState("ordersLength", []byte("0"))
	fmt.Print("PutState: ")
	fmt.Println(err)

	err = stub.PutState("usersLength", []byte("0"))
	fmt.Print("PutState: ")
	fmt.Println(err)

	return nil, nil
}

/***************************Invoke is our entry point to invoke a chaincode function**************************/

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "addProduct" {
		return t.addProduct(stub, args)
	} else if function == "addOrder" {
		return t.addOrder(stub, args)
	} else if function == "addUser" {
		return t.addUser(stub, args)
	} else if function == "setProvision" {
		return t.setProvision(stub, args)
	} else if function == "majProduct" {
		return t.majProduct(stub, args)
	} else if function == "setTrackingID" {
		return t.setTrackingID(stub, args)
	} else if function == "setState" {
		return t.setState(stub, args)
	} else if function == "setTransport" {
		return t.setTransport(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

/*********************************Query is our entry point for queries*********************************/

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

/*********************************Read and return anything from state by key*****************************************/
//args[0] : the key to use to retrieve data

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key := args[0]
	valAsbytes, err := stub.GetState(key)
	fmt.Print("valAsBytes: ")
	fmt.Println(valAsbytes)
	fmt.Println(err)

	return valAsbytes, nil
}

/***********************************Read and return a product by ref*****************************************/
//args[0] : the ref of the wanted producut

func (t *SimpleChaincode) getProductByRef(stub shim.ChaincodeStubInterface, args []string) ([]byte, int, error) {
	if len(args) != 1 {
		return nil, -1, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	var i int
	var product Product

	ref := args[0]
	fmt.Print("ref: ")
	fmt.Println(ref)

	productsLengthAsbytes, err := stub.GetState("productsLength")
	fmt.Print("ProductsLengthAsBytes: ")
	fmt.Println(productsLengthAsbytes)
	fmt.Println(err)

	length, err := strconv.Atoi(string(productsLengthAsbytes))
	for i = 0; i < length; i++ {
		productAsBytes, err := stub.GetState("product" + strconv.Itoa(i))
		fmt.Print("ProductsAsBytes " + strconv.Itoa(i) + ": ")
		fmt.Println(productsLengthAsbytes)
		fmt.Println(err)

		err = json.Unmarshal(productAsBytes, &product)
		fmt.Print("Product: ")
		fmt.Println(product)
		fmt.Println(err)

		if product.Ref == ref {
			return productAsBytes, i, nil
		}
	}

	return nil, -1, errors.New("Product not found for ref: " + ref)
}

/********************************Read and return an order by ref*************************************/
//args[0] : the ref of the wanted order

func (t *SimpleChaincode) getOrderByRef(stub shim.ChaincodeStubInterface, args []string) ([]byte, int, error) {
	if len(args) != 1 {
		return nil, -1, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	var i int
	var order Order

	ref := args[0]
	fmt.Print("ref: ")
	fmt.Println(ref)

	ordersLengthAsbytes, err := stub.GetState("ordersLength")
	fmt.Print("ordersLengthAsBytes: ")
	fmt.Println(ordersLengthAsbytes)
	fmt.Println(err)

	length, err := strconv.Atoi(string(ordersLengthAsbytes))
	for i = 0; i < length; i++ {
		orderAsBytes, err := stub.GetState("order" + strconv.Itoa(i))
		fmt.Print("ordersAsBytes " + strconv.Itoa(i) + ": ")
		fmt.Println(ordersLengthAsbytes)
		fmt.Println(err)

		err = json.Unmarshal(orderAsBytes, &order)
		fmt.Print("order: ")
		fmt.Println(order)
		fmt.Println(err)

		if order.Ref == ref {
			return orderAsBytes, i, nil
		}
	}

	return nil, -1, errors.New("order not found for ref: " + ref)
}

/******************************Add a product to the state**************************************/
//args[0] : ref of the product
//args[1] : description of the product
//args[2] : price of the product
//args[3] : quantity of the product in stock
//args[4] : critical quantity of the product in stock

func (t *SimpleChaincode) addProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	var err error
	var product Product

	fmt.Println("args[0] : " + args[0])
	fmt.Println("args[1] : " + args[1])
	fmt.Println("args[2] : " + args[2])
	fmt.Println("args[3] : " + args[3])
	fmt.Println("args[4] : " + args[4])

	productsLengthAsBytes, err := stub.GetState("productsLength")
	productsLength := string(productsLengthAsBytes)
	fmt.Print("productsLength:")
	fmt.Println(productsLength)
	fmt.Println(err)

	product.Ref = args[0]
	product.Description = args[1]
	product.Price, err = strconv.ParseFloat(args[2], 64)
	product.Quantity, err = strconv.Atoi(args[3])
	product.Critical, err = strconv.Atoi(args[4])
	product.Provision = 0

	productAsBytes, err := json.Marshal(product)
	fmt.Print("productAsBytes: ")
	fmt.Println(err)

	err = stub.PutState("product"+productsLength, productAsBytes)
	fmt.Print("PutState: ")
	fmt.Println(err)

	count, err := strconv.Atoi(productsLength)
	count++

	err = stub.PutState("productsLength", []byte(strconv.Itoa(count)))
	fmt.Print("PutState: ")
	fmt.Println(err)

	return nil, nil
}

/*********************************Set a provisioning rule to a product*********************************/
//args[0] : ref of the product
//args[1] : provisioning number

func (t *SimpleChaincode) setProvision(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	var value string
	var arguments []string
	var product Product

	fmt.Println("args[0] : " + args[0])
	fmt.Println("args[1] : " + args[1])

	arguments = append(arguments, args[0])
	value = args[1]

	productAsBytes, index, err := t.getProductByRef(stub, arguments)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Print("productAsBytes: ")
	fmt.Println(index)
	fmt.Println(err)

	err = json.Unmarshal(productAsBytes, &product)
	fmt.Print("product: ")
	fmt.Println(product)
	fmt.Println(err)

	product.Provision, err = strconv.Atoi(value)
	fmt.Print("product: ")
	fmt.Println(product)
	fmt.Println(err)

	productAsBytes, err = json.Marshal(product)
	fmt.Print("productAsBytes: ")
	fmt.Println(productAsBytes)
	fmt.Println(err)

	err = stub.PutState("product"+strconv.Itoa(index), productAsBytes)
	fmt.Print("PutState: ")
	fmt.Println(err)

	return nil, nil
}

/************Maintain stock state and deliver events when critical point is reached*********************/
//args[0] : product array
//args[1] : quantity array
//args[2] : ref of the order

func (t *SimpleChaincode) majProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	var i int
	var err error
	var orderRef string
	var product Product
	var arguments []string
	var productArray []Product
	var quantityArray []int

	fmt.Println("args[0] : " + args[0])
	fmt.Println("args[1] : " + args[1])
	fmt.Println("args[2] : " + args[2])

	orderRef = args[2]

	err = json.Unmarshal([]byte(args[0]), &productArray)
	fmt.Print("productArray:")
	fmt.Println(productArray)
	fmt.Println(err)

	err = json.Unmarshal([]byte(args[1]), &quantityArray)
	fmt.Print("quantityArray:")
	fmt.Println(quantityArray)
	fmt.Println(err)

	ordersLengthAsBytes, err := stub.GetState("ordersLength")
	ordersLength, err := strconv.Atoi(string(ordersLengthAsBytes))
	fmt.Print("ordersLength:")
	fmt.Println(ordersLength)
	fmt.Println(err)

	for i = 0; i < len(productArray); i++ {
		arguments = append(arguments, productArray[i].Ref)
		productAsBytes, index, err := t.getProductByRef(stub, arguments)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		fmt.Print("productAsBytes: ")
		fmt.Println(index)
		fmt.Println(err)

		err = json.Unmarshal(productAsBytes, &product)
		fmt.Print("product: ")
		fmt.Println(product)
		fmt.Println(err)

		if product.Quantity > quantityArray[i] {
			product.Quantity = product.Quantity - quantityArray[i]
			fmt.Print("qtfinal: ")
			fmt.Println(product.Quantity)

			productAsBytes, err = json.Marshal(product)
			fmt.Print("productAsBytes: ")
			fmt.Println(err)

			err = stub.PutState("product"+strconv.Itoa(index), productAsBytes)
			fmt.Print("PutState: ")
			fmt.Println(err)

			if product.Critical > product.Quantity {
				fmt.Println("Event : commande en cours sur le produit X")
				var customEvent = "{eventType: 'provisioningOrder', productRef:" + product.Ref + ", quantity:" + strconv.Itoa(product.Provision) + "}"
				err = stub.SetEvent("evtSender", []byte(customEvent))
				fmt.Print("Event: ")
				fmt.Println(err)
			}

		} else {
			return nil, errors.New("Insufficient stock")
		}

		arguments = nil
	}

	arguments = append(arguments, "2", orderRef)
	return t.setState(stub, arguments)
}

/********************************Add an order to the state********************************/
//args[0] : user login
//args[1] : product array
//args[2] : quantity array
//args[3] : total price of the order
//args[4] : ref of the order

func (t *SimpleChaincode) addOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}
	fmt.Println("args[0] : " + args[0])
	fmt.Println("args[1] : " + args[1])
	fmt.Println("args[2] : " + args[2])
	fmt.Println("args[3] : " + args[3])
	fmt.Println("args[4] : " + args[4])

	var err error
	var order Order
	var collis Collis

	collis.Dimension = ""
	collis.Poids = -1

	userHashAsBytes, err := stub.GetState(args[0])
	fmt.Println("userHashAsBytes:")
	fmt.Println(err)

	ordersLengthAsBytes, err := stub.GetState("ordersLength")
	fmt.Println("ordersLenghtAsBytes:")
	fmt.Println(err)

	ordersLength := string(ordersLengthAsBytes)
	fmt.Println("ordersLength:")
	fmt.Println(ordersLength)

	count, err := strconv.Atoi(ordersLength)
	fmt.Println("currenCount:")
	fmt.Println(count)

	count++
	fmt.Println("incrementCount:")
	fmt.Println(count)

	order.Ref = args[4]
	order.ClientHash = string(userHashAsBytes)
	err = json.Unmarshal([]byte(args[1]), &order.Products)
	fmt.Println("order.Products:")
	fmt.Println(order.Products)
	fmt.Println("err unmarshal args[1]:")
	fmt.Println(err)

	err = json.Unmarshal([]byte(args[2]), &order.Quantities)
	fmt.Println("order.Quantities:")
	fmt.Println(order.Quantities)
	fmt.Println("err unmarshal args[2]:")
	fmt.Println(err)

	order.TotalPrice, err = strconv.ParseFloat(args[3], 64)
	order.Collis = collis
	order.TrackingID = ""
	order.State = 1

	ordersAsBytes, err := json.Marshal(order)
	fmt.Println("err marshal ordersAsBytes:")
	fmt.Println(err)

	err = stub.PutState("order"+ordersLength, ordersAsBytes)
	fmt.Println("key:")
	fmt.Println("order" + ordersLength)
	fmt.Println("err putting state:")
	fmt.Println(err)

	err = stub.PutState("ordersLength", []byte(strconv.Itoa(count)))
	fmt.Print("PutState: ")
	fmt.Println(err)

	return nil, nil
}

/*****************************Set a trackingID to an order*****************************************/
//args[0] : trackingID value
//args[1] : ref of the order

func (t *SimpleChaincode) setTrackingID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fmt.Println("args[0] : " + args[0])
	fmt.Println("args[1] : " + args[1])

	var trackingID string
	var arguments []string
	var order Order

	trackingID = args[0]
	arguments = append(arguments, args[1])

	orderAsBytes, index, err := t.getOrderByRef(stub, arguments)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Print("orderAsBytes: ")
	fmt.Println(orderAsBytes)
	fmt.Println(index)
	fmt.Println(err)

	err = json.Unmarshal(orderAsBytes, &order)
	fmt.Print("order: ")
	fmt.Println(order)
	fmt.Println(err)

	order.TrackingID = trackingID
	order.State = 3
	fmt.Print("modifiedOrder: ")
	fmt.Println(order)

	orderAsBytes, err = json.Marshal(order)
	fmt.Print("orderAsBytes: ")
	fmt.Println(orderAsBytes)
	fmt.Println(err)

	err = stub.PutState("order"+strconv.Itoa(index), orderAsBytes)
	fmt.Print("PutState: ")
	fmt.Println(err)

	return nil, nil
}

/*******************************Set a transport infos of an order************************************/
//args[0] : collis infos
//args[1] : ref of the order
//args[3] : key of the carrier

func (t *SimpleChaincode) setTransport(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	fmt.Println("args[0] : " + args[0])
	fmt.Println("args[1] : " + args[1])
	fmt.Println("args[2] : " + args[2])

	var collis Collis
	var carrierKey string
	var arguments []string
	var order Order

	arguments = append(arguments, args[1])
	carrierKey = args[2]

	err := json.Unmarshal([]byte(args[0]), &collis)
	fmt.Print("collis: ")
	fmt.Println(collis)
	fmt.Println(err)

	userHashAsBytes, err := stub.GetState(carrierKey)
	fmt.Print("userHashAsBytes: ")
	fmt.Println(err)

	orderAsBytes, index, err := t.getOrderByRef(stub, arguments)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Print("orderAsBytes: ")
	fmt.Println(orderAsBytes)
	fmt.Println(index)
	fmt.Println(err)

	err = json.Unmarshal(orderAsBytes, &order)
	fmt.Print("order: ")
	fmt.Println(order)
	fmt.Println(err)

	order.Collis = collis
	order.CarrierHash = string(userHashAsBytes)
	fmt.Print("modifiedOrder: ")
	fmt.Println(order)

	orderAsBytes, err = json.Marshal(order)
	fmt.Print("orderAsBytes: ")
	fmt.Println(orderAsBytes)
	fmt.Println(err)

	err = stub.PutState("order"+strconv.Itoa(index), orderAsBytes)
	fmt.Print("PutState: ")
	fmt.Println(err)

	return nil, nil
}

/************************Change the statut of an order****************************************/
//args[0] : state value
//args[1] : ref of the order

func (t *SimpleChaincode) setState(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fmt.Println("args[0] : " + args[0])
	fmt.Println("args[1] : " + args[1])

	var state int
	var arguments []string
	var order Order

	state, err := strconv.Atoi(args[0])
	arguments = append(arguments, args[1])

	orderAsBytes, index, err := t.getOrderByRef(stub, arguments)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Print("orderAsBytes: ")
	fmt.Println(orderAsBytes)
	fmt.Println(index)
	fmt.Println(err)

	err = json.Unmarshal(orderAsBytes, &order)
	fmt.Print("order: ")
	fmt.Println(order)
	fmt.Println(err)

	order.State = state
	fmt.Print("modifiedOrder: ")
	fmt.Println(order)

	orderAsBytes, err = json.Marshal(order)
	fmt.Print("orderAsBytes: ")
	fmt.Println(orderAsBytes)
	fmt.Println(err)

	err = stub.PutState("order"+strconv.Itoa(index), orderAsBytes)
	fmt.Print("PutState: ")
	fmt.Println(err)

	return nil, nil
}

/*************************Add a user to the state*************************************/
//args[0] : user login
//args[1] : user password
//args[2] : user hash

func (t *SimpleChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	fmt.Println("args[0] : " + args[0])
	fmt.Println("args[1] : " + args[1])
	fmt.Println("args[2] : " + args[2])

	var err error
	var userLogin, userPassword, userHash string

	userLogin = args[0]
	userPassword = args[1]
	userHash = args[2]

	usersLengthAsBytes, err := stub.GetState("usersLength")
	fmt.Print("userLengthAsBytes: ")
	fmt.Println(err)

	err = stub.PutState(userLogin+"@"+userPassword, []byte(string(userHash)))
	fmt.Print("PutState: ")
	fmt.Println(err)

	usersLength := string(usersLengthAsBytes)
	count, err := strconv.Atoi(usersLength)
	count++

	err = stub.PutState("usersLength", []byte(strconv.Itoa(count)))
	fmt.Print("PutState: ")
	fmt.Println(err)

	return nil, nil
}