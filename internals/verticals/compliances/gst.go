package compliances

type GSTModel struct {
	SellerGSTIN       string   `json:"sellerGstin" csv:"Seller GSTIN"`
	BaseData          string   `json:"baseData" csv:"Base data"`
	Type              string   `json:"type" csv:"Type"`
	ECommerceOperator string   `json:"eCommerceOperator" csv:"e-Commerce operator"` // Handles #REF! as string
	ECommerceGSTIN    string   `json:"eCommerceGstin" csv:"e-Commerce GSTIN"`       // Handles #REF! as string
	InvoiceNo         string   `json:"invoiceNo" csv:"Invoice no"`
	InvoiceDate       string   `json:"invoiceDate" csv:"Invoice Date"` // Store as string, parse later if needed
	Quantity          int      `json:"quantity" csv:"Qty"`             // Assuming whole number quantities
	HSN               string   `json:"hsn" csv:"HSN"`
	State             string   `json:"state" csv:"State"`
	POS               string   `json:"pos" csv:"POS"` // Handles #ERROR! as string
	CustomerGSTIN     string   `json:"customerGstin" csv:"Customer GSTIN"`
	InvoiceValue      *float64 `json:"invoiceValue" csv:"Invoice Value"`    // Pointer for potential null/parse errors
	TaxableAmount     *float64 `json:"taxableAmount" csv:"Taxable amount"`  // Pointer for potential null/parse errors
	CGSTRate          *float64 `json:"cgstRate" csv:"CGST Rate"`            // Pointer for potential null/parse errors
	SGSTRate          *float64 `json:"sgstRate" csv:"SGST Rate"`            // Pointer for potential null/parse errors
	IGSTRate          *float64 `json:"igstRate" csv:"IGST Rate"`            // Pointer for potential null/parse errors
	TCSIgst           *float64 `json:"tcsIgst,omitempty" csv:"TCS IGST"`    // Pointer for potential null/empty
	TCSCgst           *float64 `json:"tcsCgst,omitempty" csv:"TCS CGST"`    // Pointer for potential null/empty
	TCSSgst           *float64 `json:"tcsSgst,omitempty" csv:"TCS SGST"`    // Pointer for potential null/empty
	TCSRate           *float64 `json:"tcsRate,omitempty" csv:"Rate_1"`      // TCS Rate (second occurrence) - Renamed header for csv tag. Pointer for potential null/empty
	Classification    string   `json:"classification" csv:"Classification"` // Handles #REF! as string
	TaxCategory       string   `json:"taxCategory" csv:"Tax category"`
	TaxType           string   `json:"taxType" csv:"Tax Type"`       // Handles #ERROR! as string
	CGSTAmount        *float64 `json:"cgstAmount" csv:"CGST Amount"` // Pointer for potential null/parse errors
	SGSTAmount        *float64 `json:"sgstAmount" csv:"SGST Amount"` // Pointer for potential null/parse errors
	IGSTAmount        *float64 `json:"igstAmount" csv:"IGST Amount"` // Pointer for potential null/parse errors
	// Rate              *float64 `json:"rate" csv:"Rate"`              // Tax Rate (first occurrence) - Pointer for potential null/parse errors
	ShippingToState   string `json:"shippingToState" csv:"Shipping to state"`
	ShippingFromState string `json:"shippingFromState" csv:"Shipping from state"`
}
