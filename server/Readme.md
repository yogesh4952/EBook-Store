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
