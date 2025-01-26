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
| **Scenario 14: Admin Access to Profile with `user:write:self` Scope**               | `PUT (defineRole(["admin"]), defineScope("user:write:self")) /api/v1/profiles/:id` | `PUT /api/v1/profiles/1` | `admin`  | `user:write:self` | `admin1`            | `admin1`              | `200 OK`              | **Admin** can update any user's profile with `self` scope.                                      |
| **Scenario 14: Admin Access to Profile with `user:write:self` Scope**               | `PUT (defineRole(["admin"]), defineScope("user:write:self")) /api/v1/profiles/:id` | `PUT /api/v1/profiles/1` | `admin`  | `user:write:self` | `admin1`            | `user1`               | `403 Forbidden`       | **Admin** cannot update another user's profile with `self` scope.                               |
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


# Permission Handling Priority

When handling permissions in an application, especially in systems with roles and scopes (like your role-based access control system), the **priority of permissions** is critical for ensuring proper access control. Here’s a structured approach to thinking about the priority of permissions:

## 1. Role-Based Permissions Over Scope-Based Permissions
   - **Roles** often provide a broad set of permissions that may include many actions, such as `admin`, `user`, `manager`, etc. These roles should generally have **higher priority** because they usually come with a large set of actions that define what the user can do.
   - **Scopes** are often more specific actions, like `"user:read:self"` or `"admin:read:all"`. While scopes offer fine-grained access control, they are **usually subordinate** to roles since roles are designed to encapsulate broad permissions.

   ### Example:
   - A user with an **admin role** (e.g., `"admin"`) could have access to every resource regardless of the specific scope, meaning `admin:read:all`, `admin:write:all`, etc.
   - A user with the role **user** (e.g., `"user"`) would have more limited access, like `user:read:self`, but cannot access the entire system as an admin.

   ### Priority:
   - **Admin > User** roles.
   - Specific scopes like `read` or `write` (based on role) are evaluated afterward.

---

## 2. Allowing or Denying Access Based on Roles
   When a request is made, first check the **role-based permissions** to determine if the user should have access. Only then should you verify the **scope-based permissions** for more fine-grained control.

   ### Flow:
   1. **Role Check**: Does the user have the required role for the resource (e.g., `admin` or `user`)?
   2. **Scope Check**: Once the role is confirmed, check if the user has the required **specific action** (e.g., `read`, `write`) and the **resource** they are trying to access (e.g., `"self"` for personal data or `"all"` for global data).

   ### Example:
   - Admin users should have **access to everything** (e.g., `admin:read:all`), but a **user** role with `"user:read:self"` can only read their own data.

   ### Priority:
   - **Role Check**: Determines if the user is **allowed** to access any data.
   - **Scope Check**: Further **restricts** access based on specific resources/actions.

---

## 3. Conflicting Permissions
   In scenarios where conflicting permissions might arise, you need to define a clear rule for how conflicts are resolved:

   ### Priorities to Consider:
   - **Deny by Default**: This is a common principle in access control. If there’s any ambiguity or conflict in permissions, **deny access** by default and only allow explicit permissions.
   - **Most Restrictive Permission**: If a user has both `admin:read:all` and `user:read:self`, prioritize the **more restrictive** scope in some cases. For instance, an admin might have access to more data, but if there’s a resource-specific permission tied to `"self"`, that might take precedence to ensure users can’t access others’ data.

---

## 4. Resource Ownership and Self-Access
   **Self-access** permissions generally have higher importance when accessing personal data (like profiles) because it ensures privacy and compliance with data protection rules.

   ### Priority:
   - If the scope is `"self"`, it usually has a **higher priority** for user privacy, meaning users can only access their own data.
   - If an **admin** tries to access a user’s personal data (e.g., via a `"user:read:self"` scope), the system should ideally deny access unless the resource is explicitly allowed for admin access.

---

## 5. Explicit Permissions over Implicit Permissions
   If a user’s permissions are **explicitly granted** (e.g., through a role or scope), these permissions should always **override** any default or general permissions.

   ### Example:
   - An admin role might give full access, but if you have an explicit scope check (like `"user:write:self"`), you still need to ensure that the action is checked.

   ### Priority:
   - **Explicit permissions** take precedence over **implicit permissions** or default rules.

---

## Prioritization Strategy for Permissions Handling

1. **Role-Based Access Control (RBAC)** (Primary Check):
   - **Admin roles** generally have higher access than **user roles**.
   - **Check the user’s role** to identify broad access rights (e.g., `"admin:read:all"`, `"admin:write:all"`).
   - If the user is an admin, they likely have access to most resources.

2. **Scope-Based Access Control** (Secondary Check):
   - **Scalable and fine-grained permissions** that apply to specific actions (e.g., `"user:read:self"`, `"user:write:self"`).
   - **Self-scopes** (for personal data) should be considered higher priority for privacy and data protection.

3. **Permission Deny (Fallback Rule)**:
   - If no matching permission is found or if there’s an ambiguity in the roles/scopes, the system should **deny access by default** (a best practice in security).

---

## Example of a Permission Priority Flow

- **GET /accounts (Role: Admin)**:
   - Check if the role is `"admin"`. If true, **allow access**.
   - If not `"admin"`, check scope (`user:read:self`, etc.).

- **GET /profiles/:id (Role: User, Scope: user:read:self)**:
   - Check if the user has `"user"` role. If yes, check if they’re trying to access their own profile.
   - If the scope is `"self"`, restrict access to the **authenticated user's own profile only**.

- **GET /transactions (Role: Admin, Scope: admin:read:all)**:
   - Admin role allows access, irrespective of scope restrictions.

- **GET /transactions (Role: User, Scope: user:read:self)**:
   - User can only read their own transactions. Admin cannot access a user's transactions without explicit permission.

---

## Final Priority Order
1. **Role-Based Access Control (RBAC)** (admin > user > others)
2. **Scope-Based Access Control** (read/write permissions for self vs. all)
3. **Specific Resource Access** (self vs all) for personal data, with **privacy protections** prioritized.
4. **Permission Denial (Fallback Rule)** if no valid permission is found.



