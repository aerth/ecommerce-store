Pull Requests are welcome. Work in progress.

## Make a one time payment, Add user to Database.

Should be a simple task.

Here is the flow:

1. Someone clicks "Purchase", they get a form for their New username and password. (for our site)
2. They get a link to use paypal. (near future: credit card form)
3. Click paypal, OFFSITE LINK to paypal.com login form, choose payment method
4. RETURN to our site, now with a payment ID and Payer ID.
5. Allow user to "confirm" before the real purchase execution
6. Server sends the "confirmation" to paypal, we get "Approval State: approved"
7. Email / password gets saved (plaintext for testing) in sub.db
