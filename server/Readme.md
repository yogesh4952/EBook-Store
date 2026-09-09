# PAYMENT FLOW

usercheckout -> backend calls provider -> backend_get_link -> redirected to client -> user fill details -> open a webhook in backend -> succes= redirect to succes and verify in backend

two level security

frontend and backend both verifies

#### /webhook/payment

- verifies cryptographic signature
- change the status to paid

### /verify-payment?id=lfjakjfs

- frontend redirected to success page
- then call this endpoint
- check whether the status="PAID" or not
- based upon it redirect and show message

The OrderItem repo method itself is fine
CreateInBatches(data) correctly inserts rows into the order_items table. Nothing wrong with that line.
The real problems

1. It's never called — and not on the interface. IOrderRepo (order_repo.go:10-12) only declares PlaceOrder. The OrderItem method on line 27 exists but is dead code; the service can't call it through the interface.
2. The service never builds order items. In order_service.go, PlaceOrder computes total_price from items (line 55-63) but discards them — it never constructs []\*models.OrderItem nor inserts them. So the order_items table stays empty even after a successful order.
3. OrderId would be empty. Each OrderItem needs OrderId filled to associate with the parent order. But data.ID is only populated after PlaceOrder creates the order (line 73), so you must build items after that call, using data.ID.
   How to debug/fix
   Add OrderItem to the interface and call it in the service after the order is created:
   // interface
   type IOrderRepo interface {
   PlaceOrder(ctx context.Context, data *models.Order) error
   OrderItem(ctx context.Context, data []*models.OrderItem) error
   }
   // in service, after PlaceOrder succeeds
   var items []\*models.OrderItem
   for _, val := range orderPayload.Items {
   book, _ := serv.bookRepo.FindById(ctx, val.BookId)
   items = append(items, &models.OrderItem{
   OrderId: data.ID, // now populated
   BookId: val.BookId,
   Quantity: val.Quantity,
   UnitPrice: uint(book.Price),
   })
   }
   if err := serv.orderRepo.OrderItem(ctx, items); err != nil {
   return nil, err
   }
   Note: UnitPrice here is uint, but book.Price is likely float32 — you'll need to handle that conversion appropriately (and fix the res.Items subtotal/price mismatches too, since UnitPrice as uint can't hold decimals).
   ▣ Build · Big Pickle · 15.0s
