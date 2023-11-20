# email-controller

Using this project as your skeleton, write a kubernetes controller that will send an email. Then, write a script that process the provided `contacts.csv` file and utilizes the API & controller you've built to send an email to each contact.

You will need to design your API (see `./api/v1alpha1/emailrequest_types.go`) and implement the `EmailRequest` controller (`./internal/controller/emailrequest_controller.go`). 

To each contact, send an email with a subject of `Hello` and a body of `Hello, {name}!`. 

A fake email service has been provided for you to use in your controller (`./pkg/fakeemail`). 

We expect that you will spend roughly 2 hours on this assignment, but please spend no more than 4 hours. An incomplete submission followed-up with good answers in the first interview round is better than a perfect submission and poor answers during the interview.

## Requirements
- Each email must be sent at most once
- EmailRequests should have configurable retry behavior that applies to retriable errors. This must be configurable at the API level, not controller-wide.
- Emails that are "bounced" should be treated as permanent failures and not retried.
- Emails that are "blocked" should be treated as temporary failures and retried based on the configured retry behavior.
- Invalid email addresses should be treated as permanent failures and not retried.
- Your status API for `EmailRequest` should include status conditions that indicate whether the email has been sent successfully (among any other status information you deem important to expose)

## Bonus points

Submissions that go "above and beyond" might include any or all of the following. Do note that this is optional, but encouraged, based on your available time.
- Regex-based validation of email addresses at the API level. This could be via kubebuilder markers or via a validating webhook. You can use the same regex that the fakeemail service uses: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$` 
- Tests that cover your controller's retry semantics and status API. See the [kubebuilder docs](https://book.kubebuilder.io/cronjob-tutorial/writing-tests) on writing tests for guidance.


## Next steps

When you have completed the assignment, please return your full implementation (the entire repo) in a zip file. A follow-up interview will then be scheduled; in this interview, among other things, we will discuss the decisions you made in your submission and what motivated these decisions.


--------------------


## Running the controller 
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -k config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/email-controller:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/email-controller:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

