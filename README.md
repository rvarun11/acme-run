# ACME RUN - Liuyin, Samkith, Varun

Welcome to our implementation of ACME RUN, a collaborative project developed by Liuyin, Samkith, and Varun as part of CAS 735 (Fall 2023).

## Installation

1. Clone the Repository:
```bash
git clone https://github.com/CAS735-F23/macrun-teamvsl.git
cd macrun-teamvsl
```

2. Run Docker Compose
```bash
docker compose up
```

## Getting Started

After successfully setting up the project, follow these steps to get started:

1. **Import Postman Collection**:
- Import the provided Postman Collection from the repository into your Postman workspace or use the [link](https://winter-satellite-393249.postman.co/workspace/cas-735~2906f288-5f3e-4839-8f70-f7f36274cd09/collection/14312203-b6260f24-54b8-4d85-8684-dcd9821a3545?action=share&creator=14312203).
- The collection includes predefined API requests for ACME RUN.
- The collection is divided into 9 folders, the first for initializing the application and the others for the different scenarios. Each folder has its description explaining what it does.
- All the folders and APIs are ordered to be run from top to bottom. We recommend you to each API one by one while simultaneously seeing the logs in the terminal.

2. **Run APIs**:
- Services use ports 8010 to 8014 and 5432 (for PostgreSQL) by default; ensure they are available for API functionality.
- Execute the API requests from the top of the Postman Collection.
- Explore and interact with the various endpoints available in ACME RUN.

## Testing

Run the tests with the following command:
```bash
./run_tests.sh
```
If you don't have Golang configured locally, you can use the preconfigured devcontainer. Simply reopen the repository in the Devcontainer (or Codespaces) and run the script. For additional information about the tests, refer to the [test description here](#TESTS.md).

## Contributors

- Liuyin Shi (shil9@mcmaster.ca)
- Samkith Jain (kishors@mcmaster.ca)
- Varun Rajput (rajpuv2@mcmaster.ca)
