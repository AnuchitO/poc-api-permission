### **Test Case Scenarios for Role-Based and Scope-Based Access Control**

| **Test Case Name**                                                                  | **Endpoint Definition**                                                            | **Client Request**       | **Role** | **Scope**         | **User ID (Claim)** | **Resource Owner ID** | **Expected Response** | **Response Code**                                                                               | **Explanation** |
| ----------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ------------------------ | -------- | ----------------- | ------------------- | --------------------- | --------------------- | ----------------------------------------------------------------------------------------------- | --------------- |
| **Scenario 1: User Access to Admin-Only Endpoint**                                  | `GET (defineRole(["admin"]) /api/v1/accounts`                                      | `GET /api/v1/accounts`   | `user`   | N/A               | `user1`             | N/A                   | `403 Forbidden`       | **User** does not have access to an `admin`-only endpoint.                                      |
| **Scenario 2: Admin Access to Admin-Only Endpoint**                                 | `GET (defineRole(["admin"]) /api/v1/accounts`                                      | `GET /api/v1/accounts`   | `admin`  | N/A               | `admin1`            | N/A                   | `200 OK`              | **Admin** can access the `admin`-only endpoint.                                                 |
| **Scenario 3: User Access to Admin-Only Profile Endpoint**                          | `GET (defineRole(["admin"]) /api/v1/profiles`                                      | `GET /api/v1/profiles`   | `user`   | N/A               | `user1`             | `user1`               | `403 Forbidden`       | **User** cannot access `admin`-only profile endpoint.                                           |
| **Scenario 4: Admin Access to Admin-Only Profile Endpoint**                         | `GET (defineRole(["admin"]) /api/v1/profiles`                                      | `GET /api/v1/profiles`   | `admin`  | N/A               | `admin1`            | `user1`               | `200 OK`              | **Admin** can access any profile due to `admin` role.                                           |
| **Scenario 5: Admin Access to Account with `admin:read:all` Scope**                 | `GET (defineRole(["admin"]), defineScope("admin:read:all")) /api/v1/accounts/:id`  | `GET /api/v1/accounts/2` | `admin`  | `admin:read:all`  | `admin1`            | `user2`               | `200 OK`              | **Admin** with `admin:read:all` scope can access any account.                                   |
| **Scenario 6: User Access to Account with `user:read:self` Scope**                  | `GET (defineRole(["user"]), defineScope("user:read:self")) /api/v1/accounts/:id`   | `GET /api/v1/accounts/1` | `user`   | `user:read:self`  | `user1`             | `user1`               | `200 OK`              | **User** can access their own account with `self` scope.                                        |
| **Scenario 7: User Access to Another User's Account with `user:read:self` Scope**   | `GET (defineRole(["user"]), defineScope("user:read:self")) /api/v1/accounts/:id`   | `GET /api/v1/accounts/2` | `user`   | `user:read:self`  | `user1`             | `user2`               | `403 Forbidden`       | **User** cannot access another user's account with `self` scope.                                |
| **Scenario 8: Admin Access to Account with `user:read:self` Scope**                 | `GET (defineRole(["admin"]), defineScope("user:read:self")) /api/v1/accounts/:id`  | `GET /api/v1/accounts/1` | `admin`  | `user:read:self`  | `admin1`            | `user1`               | `200 OK`              | **Admin** can access any account with `self` scope for the user (admin has broader permission). |
| **Scenario 9: User Access to Profile with `user:read:self` Scope**                  | `GET (defineRole(["user"]), defineScope("user:read:self")) /api/v1/profiles/:id`   | `GET /api/v1/profiles/1` | `user`   | `user:read:self`  | `user1`             | `user1`               | `200 OK`              | **User** can access their own profile with `self` scope.                                        |
| **Scenario 10: User Access to Another User's Profile with `user:read:self` Scope**  | `GET (defineRole(["user"]), defineScope("user:read:self")) /api/v1/profiles/:id`   | `GET /api/v1/profiles/2` | `user`   | `user:read:self`  | `user1`             | `user2`               | `403 Forbidden`       | **User** cannot access another user's profile with `self` scope.                                |
| **Scenario 11: Admin Access to Profile with `user:read:self` Scope**                | `GET (defineRole(["admin"]), defineScope("user:read:self")) /api/v1/profiles/:id`  | `GET /api/v1/profiles/1` | `admin`  | `user:read:self`  | `admin1`            | `user1`               | `200 OK`              | **Admin** can access their own profile with `self` scope.                                       |
| **Scenario 12: User Access to Profile with `user:write:self` Scope**                | `PUT (defineRole(["user"]), defineScope("user:write:self")) /api/v1/profiles/:id`  | `PUT /api/v1/profiles/1` | `user`   | `user:write:self` | `user1`             | `user1`               | `200 OK`              | **User** can update their own profile with `self` scope.                                        |
| **Scenario 13: User Access to Another User's Profile with `user:write:self` Scope** | `PUT (defineRole(["user"]), defineScope("user:write:self")) /api/v1/profiles/:id`  | `PUT /api/v1/profiles/2` | `user`   | `user:write:self` | `user1`             | `user2`               | `403 Forbidden`       | **User** cannot update another user's profile with `self` scope.                                |
| **Scenario 14: Admin Access to Profile with `user:write:self` Scope**               | `PUT (defineRole(["admin"]), defineScope("user:write:self")) /api/v1/profiles/:id` | `PUT /api/v1/profiles/1` | `admin`  | `user:write:self` | `admin1`            | `user1`               | `200 OK`              | **Admin** can update any user's profile with `self` scope.                                      |
| **Scenario 15: Admin Access to Profile with `admin:write:all` Scope**               | `PUT (defineRole(["admin"]), defineScope("admin:write:all")) /api/v1/profiles/:id` | `PUT /api/v1/profiles/1` | `admin`  | `admin:write:all` | `admin1`            | `user1`               | `200 OK`              | **Admin** can update any profile with `admin:write:all` scope.                                  |

---

### **Explanation of Columns:**

- **Test Case Name**: The name or description of the scenario being tested.
- **Endpoint Definition**: The API endpoint definition, including the `defineRole` function to specify roles that can access the endpoint.
- **Client Request**: The actual HTTP request being made (e.g., `GET /api/v1/accounts`).
- **Role**: The role of the client/user making the request (e.g., `admin`, `user`).
- **Scope**: The scope specified for the request, if applicable (e.g., `user:read:self`, `admin:read:all`).
- **User ID (Claim)**: The user ID contained in the JWT claim, which can be used for access control.
- **Resource Owner ID**: The ID of the resource being requested, if relevant to the access control decision.
- **Expected Response**: The expected result of the API request based on the defined roles and scopes.
- **Response Code**: The expected HTTP response code for the given request.
- **Explanation**: Why the expected response code is the correct one for this scenario.

## Suggested Approach:

Use specific permissions like "user:read:self", "user:write:self", and "admin:read:all". This approach gives you full control and is clear.

You can reserve "self" and "all" for use as values within the scope, like:
- "user:read:{self}" — where {self} is dynamically replaced with "the current user's data" or "the current resource".
- "admin:read:{all}" — where {all} allows admin to access all resources.

### Example:
#### With Full Scopes:
- GET (defineRole(["user"]), defineScope("user:read:self")) /api/v1/profiles/:id
Meaning: User can only read their own profile (self as the resource).
GET (defineRole(["admin"]), defineScope("admin:read:all")) /api/v1/profiles/:id
Meaning: Admin can read any user's profile (all as the resource).

#### With Simplified Scopes:
- GET (defineRole(["user"]), defineScope("self")) /api/v1/profiles/:id
    - Meaning: User can only read their own profile.
- GET (defineRole(["admin"]), defineScope("all")) /api/v1/profiles/:id
    - Meaning: Admin can read any profile.

### Conclusion:
- For better flexibility, granularity, and clarity, I would recommend using specific scopes like "user:read:self" and "admin:read:all". This gives you more explicit control and ensures better clarity in defining roles and permissions.
- For simpler cases, "self" and "all" can be useful, but remember that you'll need to pair them carefully with the user’s role to avoid ambiguity.
