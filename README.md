# email-controller

Using this project as your skeleton, write a kubernetes controller that will send an email. Then, write a script that process the provided `contacts.csv` file and utilizes the API & controller you've built to send an email to each contact.

You will need to design your API (see `./api/v1alpha1/emailrequest_types.go`) and implement the `EmailRequest` controller (`./internal/controller/emailrequest_controller.go`). 

To each contact, send an email with a subject of `Hello` and a body of `Hello, {name}!`. 

A fake email service has been provided for you to use in your controller (`./pkg/fakeemail`). **Do NOT modify the fake email client.**

We expect that you will spend roughly 2 hours on this assignment, but please spend no more than 4 hours. An incomplete submission followed-up with good answers in the first interview round is better than a perfect submission and poor answers during the interview.

## Requirements
- Each email must be sent at most once, including in the face of network failures, program panics, or runtime errors.
  - You can ignore cases such as power loss, OOMKill, and other unexpected shutdown events.
- EmailRequests should have configurable exponential-backoff retry behavior that applies to retriable errors. This should be configurable on a per-EmailRequest basis rather than via controller-wide configuration. (Hint: Design the EmailRequest API well!)
  - The base delay and maximum delay components of your exponential backoff behavior must be configurable. You can ignore jitter.
- Emails that are "bounced" should be treated as permanent failures and not retried.
- Emails that are "blocked" should be treated as temporary failures and retried based on the configured retry behavior.
- Invalid email addresses should be treated as permanent failures and not retried.
- Your status API for `EmailRequest` should include [status conditions](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Condition) that indicate whether the email has been sent successfully (among any other status information you deem important to expose).

## Bonus points

Submissions that go "above and beyond" might include any or all of the following. Do note that this is optional, but encouraged, based on your available time.
- Tests that cover your controller's retry semantics and status API. See the [kubebuilder docs](https://book.kubebuilder.io/cronjob-tutorial/writing-tests) on writing tests for guidance.
- Regex-based validation of email addresses at the API level. This could be via kubebuilder markers or via a validating webhook. You can use the same regex that the fakeemail service uses: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
  - NOTE: if you implement this, your controller should still handle invalid email address errors from the fake email service, regardless of whether or not you are catching them via the same regex at the API level.


## Next steps

When you have completed the assignment, please return your full implementation (the entire repo) in a zip file. A follow-up interview will then be scheduled; in this interview, among other things, we will discuss the decisions you made in your submission and what motivated these decisions.


--------------------


## Running the controller 
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

Create a kind cluster:
```sh
kind create cluster -n influxdata-test
```

### Building the project
You can lint, compile, and perform all appropriate code-generation with the default make target. This target name is not required, but you can also refer to the target
by name with either `make all` or `make build`.

```sh
make 
```

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

You will then need to re-run `make install` to apply any generated changes to your CRDs.

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

