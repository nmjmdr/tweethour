# Tweethour

Tweethour is a RESTful API server to obtain the 'per hour frequency of tweets' for a given twiiter user for the day. Tweethour fetches the tweets for a user and classified them into 24 hours of the day (0 to 1, 1 to 2, ... 22 to 23) and computes the frequency of tweets in each hour.

To accomplish this, Tweethour uses:

  - Twitter's user [timeline API][timeline-api]
  - A http based REST server written in GO



### Version
1.0.0

### Design

Tweethour has to perform the following activites in order to compute the per-hour frequency of tweets for a given user:

* Get an Authentication token from twiiter to enable the service to get the user's timeline
* Tweethour uses [application-only authentication][app-only-auth] to authenticate itself to Twitter's API
* Once a authentication token is obtained successfuly, it needs to be attached to the request that the service invokes.
* The service invokes the user's timeline request. Note that the API will not response more than 200 tweets per request.
> **This means that if the user is a prolific tweeter and has tweeted more than 200 times on a given day, the first request to Twitter's user timeline API only returns the first 200 tweets for the day. We need to get all the tweets for the day.** 
* In order to acoomplish this:
  * Repeatedly invoke the request to get user's timeline, until the we come across a tweet that has not been created today
  * In order to obtain the next set of tweets, Twitter API provides two parameters: max_id and since_id. It is important that these parameters are set to obtain the next set of tweets. The design and code takes this into account.
  * For further details about max_id and since_id parameters reference: [Working with timelines][work-tl]
* Once all the tweets created on this day are fecthed, the program needs to classify the tweets by hour

There are three distinct steps in the above flow:

  - Get the authentication token
  - Fetch the Tweets for the user
  - Classify the tweets, (fetching more tweets - if the user has been prolific)

Tweethour has three components to accomplish this:

  - Token Generator (TokenGen)
  - User's Timeline fetcher (Timeline)
  - Histogram - to compute per-hour frequency, (Histogram)

Historgram is the main component that references and invokes the other two components. It performs the following steps:
1. Gets the authentication token (from TokenGen) and maintains it for all the future requests
2. Fetchtes the user's twwets by invoking Timeline component. It takes a call on whether to fetch the user's timeline more than once (if the user has been prolific) 
3. Computes the per-hour frequency of tweets for the day

Histogram is used by the RESTful API (server component) to fetch the per-hour frequency of tweets for a given user.

#### API Endpoints

The server exposes the following end-points:

1. GET: /histogram/{username} -- returns the JSON encoded per-hour frequenct of tweets for the given user

The format of the response for histogram is as follows:
```json
[
  {
    "from-hour": "2016-05-08T00:00:00Z",
    "to-hour": "2016-05-08T01:00:00Z",
    "frequency": 0
  },
  {
    "from-hour": "2016-05-08T01:00:00Z",
    "to-hour": "2016-05-08T02:00:00Z",
    "frequency": 1
  },
```

* "from-hour" indicates the start of the hour
* "to-hour" indicates the end of the hour
* frequency indicates the number of the tweets the user has made in that hour
* Note all times are in UTC

#### Http Response Status codes

1. Successful operation should return a response with status code 200 OK
2. If a given username (for getting the histogram) is not found, then a 404 Not found response code is returned
3. Any other errors, return appropriate HTTP status codes

Error messages are returned in the following JSON format:
```json
{
    "error-message":"Request to get user's timeline failed, 
    Status : Not Found",
    "twitter-api-response-status":"Not Found"
}
```

#### Http Methods Allowed
Only "GET" and "OPTIONS" operations are allowed on the endpoints. Invoking GET method will perform the operation and return the response as discuused above. 

OPTIONS returns the description of the methods allowed for the endpoint, which in this is only GET for all the endpoints. 

For exampple OPTIONS: /histogram/{username} returns:

```json
{
 "GET" : {
     "description" : "GET: /"
 }
}
```

#### Authentication token generation

In order to generate the authentication token, we need to encode the consumer key and consumer token for the application in a specific format prescribed by Twitter.

**The consumer key and the consumer secret has to be provided in tokengen.go file. The file contains a dummy token.**


### Setup

Clone the respository to a directory:
```sh
$ git clone https://nmjmdr1@bitbucket.org/nmjmdr1/tweethour.git tweethour
```
To compile and build, Tweethour requires: [GO](https://golang.org/) 

Once GO has been installed, execute the following commands to build the project

```sh
$ cd tweethour
$ go get github.com/gorilla/mux
(This will install the mux router, which Tweethour uses to route the http requests)
$ go build
```

This should build the executable "tweethour". In order to run the server, use the following command:


```sh
$ tweethour <port-number> 
Example:
$ tweethour 8090
Starts the tweethour API server on port 8090. The requests can be made to http://localhost:8090/
```

### Things yet to be done

* Logging API metrics
* Response caching 
  * Determine the interval of time within which if the API is invoked again for the same user, indicate the client to use the cached content  
* The consumer key and consumer secret is maintained in the source code for now. This can be moved to a encyrpted file and the program made to accept the encyrption key as a command line parameter.  A new Twitter account has been created to create the consumer key and secret for the application.




   [timeline-api]: https://dev.twitter.com/rest/reference/get/statuses/user_timeline
   [app-only-auth]: https://dev.twitter.com/oauth/application-only
   [work-tl]: https://dev.twitter.com/rest/public/timelines
   
   
