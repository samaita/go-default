This layer will act as the presenter. Decide how the data will presented. Could be as REST API, or HTML File, or gRPC whatever the delivery type. 

This layer also will accept the input from user. Sanitize the input and sent it to Usecase layer.
For my sample project, I’m using REST API as the delivery method. 

Client will call the resource endpoint over network, and the Delivery layer will get the input or request, and sent it to Usecase Layer.
This layer will depends to Usecase Layer.