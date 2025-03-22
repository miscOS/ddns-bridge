# DDNS Bridge
DDNS Bridge is a lightweight Go application that serves as an interface for updating DNS records via a simple URL. It supports providers with a web API, allowing dynamic DNS updates with minimal effort.

Features:

 - Easy-to-use HTTP endpoint for DNS updates
 - Supports multiple DNS providers with API integration
 - Secure and efficient implementation in Go

Perfect for automating DDNS updates on your own terms!

## Installation

We recommend using DDNS Bridge via Docker

```
docker-compose.yaml

```

## Configuration

User Authentication is secured by a random generated secret on startup. The environment variable `SECRET_KEY` provides an option to use your own secret that remains valid across restarts. 

If the environment variable `REGISTER_KEY` is set, it is necessary to specify the value when creating new users

### Database
By default, all data in an SQLite database at `/app/data/db.sqlite`. 

Optionally, an external MySQL (MariaDB) database can be used by setting the environment variables:
 - `DB_HOST` (MySQL host_ip:port)
 - `DB_USER` (MySQL username)
 - `DB_PASS` (MySQL password)
 - `DB_NAME` (MySQL databasename)

## API Endpoints

### Authentification
#### `POST` /api/user/signup
Creates a new user.

Parameters:
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `key`| optional | string | Value of the environment variable `REGISTER_KEY` |

Body (JSON):
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `username`| required | string | - |
| `password`| required | string | - |

#### `POST` /api/user/login
Return the token of the provided user.

Body (JSON):
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `username`| required | string | - |
| `password`| required | string | - |

### Userinformationen
#### `GET` /api/user
*Requires Token Header*  
Return information of the signed in user.

### Webhooks
#### `GET` /api/webhook/{webhook}
*Requires Token Header*  
List information of webhooks.

Parameters:
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `webhook`| optional | string | Limit the result to specified `webhook` |

#### `POST` /api/webhook
*Requires Token Header*  
Creates a new webhook for the signed in user.

Body (JSON):
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `name`| required | string | Unique name of the new `webhook` |

#### `DELETE` /api/webhook/{webhook}
*Requires Token Header*  
Removes webhook with the specified `webhook`.

Parameters:
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `webhook`| required | string | ID of the `webhook` |

### Tasks
#### `GET` /api/webhook/{webhook}/task/{task}
*Requires Token Header*  
List `task` information of the `webhook`.  

Parameters:
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `webhook`| required | string | ID of the `webhook` |
| `task`| optional | string | Limit the result to specified `task` |

#### `POST` /api/webhook/{webhook}/task
*Requires Token Header*  
Create a new `task` under the specified `webhook`

Parameters:
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `webhook`| required | string | ID of the `webhook` |

Body (JSON):
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `name`| required | string | Unique name of the new `task` |
| `service`| required | string | Service provider identifier |
| `service_params`| required | object | Service specific parameters |

#### `DELETE` /api/webhook/{webhook}/task/{task}
*Requires Token Header*  
Removes the specified `task`

Parameters:
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `webhook`| required | string | ID of the `webhook` |
| `task`| required | string | ID of the `task` |


### Update
#### `GET` /update?token={token}&ipv4={ipv4}&ipv6={ipv6}
Updates the webhook specified by its `token`

Parameters:
| Name | Required |  Type  | Description |
| ----:|:--------:|:------:|-----------------|
| `token`| required | string | Update `token` of desired the `webhook` |
| `ipv4`| optional | string | At least one of `ipv4` or `ipv6` needs to be provided |
| `ipv6`| optional | string | At least one of `ipv4` or `ipv6` needs to be provided |

## Support

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/K3K2OQ0GL)