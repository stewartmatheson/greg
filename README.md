# Greg

## The Problem
For teams developing services that consist of more than one application docker has become a vital tool. Part of the appeal of docker development is a tool called docker-compose. Docker compose allows for stacks of dockerized applications to be described in one yaml based configuration file. Many teams commit this file to a SCM to share among developers. This provides a powerful work flow as the docker compose file describes a whole stack of applications as well as instructions on downloading and running them in an automated way. The draw back to this process is that different developers require the docker-compose file to be configured in many different ways during the coarse of normal development. This becomes cumbersome and error prone
