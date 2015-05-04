var process_request = function(request, response) {
	console.log("> middleware:js:process_request");
	console.log(JSON.stringify(response));

    return;
}

var process_response = function (request, response) {
	console.log("< middleware:js:process_response");
	console.log(JSON.stringify(request));
	console.log(JSON.stringify(response));

	return response;
}
