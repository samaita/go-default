This layer will act as the business process handler. Any process will handled here. This layer will decide, which repository layer will use. And have responsibility to provide data to serve into delivery. Process the data doing calculation or anything will done here.

Usecase layer will accept any input from Delivery layer, that already sanitized, then process the input could be storing into DB , or Fetching from DB ,etc.

This Usecase layer will depends to Repository Layer