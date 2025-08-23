# About the project: A backend system to manage inventory management
This project is a RESTful API built with Golang (Gin Framework) that manages inventory and sales transactions of an e-commerce.
It includes features such as CRUD operations, authentication with JWT, role-based access control, middleware integration, and deployment to VPS.

# Tech stack
1. Language: Golang 1.24.2
2. Framework: Gin Web Framework
3. ORM: GORM
4. Database: MySQL
5. Authentication: JWT

# Database Design
The database is designed with multiple entities:
1.  Users
    Stores user information including username, email, password hash,  role (admin or user), and account status. Users can have multiple addresses.
2.  Addresses
    Represents the addresses of users. Each address belongs to one  user and includes details like receiver name, phone number, address line, city, province, and postal code.
3.  Products
    Stores product information such as name, description, price, stock, and optional expiry year.
4.  Product Images
    Each product can have multiple images. One image can be marked as the primary image.
5.  Product Categories
    Products can belong to multiple categories (many-to-many relationship). Categories store name and description.
6.  Orders
    Stores order details including user, address, status, subtotal, shipping fee, total amount, and timestamps.
7.  Order Items
    Represents products within an order, storing the product ID, quantity, and price at the time of order.
8.  Payments
    Stores payment information for orders, including invoice ID, payment type, and status. Each payment belongs to one order.
9.  Cart
    Stores shopping cart information for users. Each cart can have multiple cart items.
10. Cart Items
    Represents products added to a user's cart, including product ID and quantity.
11. Action Logs
    Logs actions performed by users or admins, including actor type, actor ID, action, and the entity affected (users, addresses, orders, payments, products, etc.).

# Installation and Setup
1. Clone the repository
2. Create a database, sql script provided in this repository
3. Create an .env (not included in this repository) file which includes:
APP_PORT=8080

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_db_password
DB_NAME=ecommerce

JWT_ACCESS_SECRET=your_jwt_access_secret
JWT_REFRESH_SECRET=your_jwt_refresh_secret
ACCESS_TTL_MIN=15
REFRESH_TTL_DAYS=7

SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your_email@example.com
SMTP_PASS=your_email_password
FROM_EMAIL=your_email@example.com
APP_ENV=dev

CLOUDINARY_URL=your_cloudinary_url
XENDIT_API_KEY=your_xendit_api_key
3. Run the server using go run main.go on the terminal

# Documentation
## Endpoints
### 1. Authentication
- `POST /auth/register` – Register a new user
- `POST /auth/login` – Login user and receive JWT tokens
- `POST /auth/verify-otp` – Verify OTP for account activation
- `POST /auth/refresh` – Refresh JWT tokens
- `POST /auth/logout` – Logout user
- `GET /auth/profile` – Get current user profile (user or admin)
- `GET /auth/admin/dashboard` – Admin-only dashboard

### 2. Users & Addresses
- `POST /user/addresses` – Create a new address (user only)
- `GET /user/addresses` – Get user's addresses
- `PATCH /user/addresses` – Update an address
- `DELETE /user/addresses/:id` – Delete an address
- `GET /admin/addresses` – Get all addresses (admin only)
- `POST /admin/addresses/:id/recover` – Recover a deleted address

### 3. Products & Categories
- `GET /products` – List all products (public)
- `GET /products/:productId` – Get product details (public)
- `POST /admin/products` – Create product (admin only)
- `PATCH /admin/products/:productId` – Update product (admin only)
- `DELETE /admin/products/:productId` – Delete product (admin only)
- `POST /admin/products/:productId/recover` – Recover product (admin only)
- `POST /admin/products/:productId/images` – Upload product image (admin only)
- `DELETE /images/:imageId` – Delete product image (admin only)
- `POST /images/:imageId/recover` – Recover deleted image (admin only)
- `GET /categories` – List all categories (public)
- `POST /admin/categories` – Create category (admin only)
- `PATCH /admin/categories/:id` – Update category (admin only)
- `DELETE /admin/categories/:id` – Delete category (admin only)
- `POST /admin/categories/:id/recover` – Recover category (admin only)

### 4. Cart
- `GET /user/cart` – Get current user's cart
- `POST /user/cart/items` – Add item to cart (user only, rate-limited)
- `DELETE /user/cart/items/:id` – Remove item from cart (user only, rate-limited)

### 5. Orders
- `POST /user/orders` – Create order from cart (user only, rate-limited)
- `GET /user/orders` – Get user's orders (user only, rate-limited)
- `GET /admin/orders` – Get all orders (admin only)
- `PUT /admin/orders/:id/status` – Update order status (admin only)

### 6. Payments
- `POST /user/payments/xendit` – Create payment via Xendit (user only, rate-limited)
- `GET /user/payments` – Get user's payments (user only, rate-limited)
- `GET /admin/payments` – Get all payments (admin only)
- `POST /admin/payments/webhook/xendit` – Xendit payment webhook (admin only)

### 7. Action Logs & Reports
- `GET /admin/logs` – Get all action logs (admin only)
- `GET /admin/logs/:id` – Get specific action log (admin only)
- `GET /admin/reports/selling` – Sales report (admin only)
- `GET /admin/reports/stock` – Stock report (admin only)

**Note:**  
Routes are protected by JWT authentication (user or admin) and logger middleware to create loggings.  
Some routes have rate limits to prevent abuse (e.g., cart actions, orders, payments).

# Deployment:
- This project is deployed on an Ubuntu VPS using Docker and automated via GitHub Actions.
- CI/CD: Deployment triggered on push to main branch using GitHub Actions workflow.
- Dockerized: The app runs inside a Docker container.
- SSH Access: GitHub Actions connects to the VPS via SSH to pull the latest code and rebuild the container.
- Port Mapping: The container exposes the application on the VPS, e.g., 8007:8000.
- Environment Variables: .env file on VPS provides all necessary configs (DB, JWT, SMTP, Xendit, Cloudinary, etc.).

**Author:**
Developed by Zahra<br>
Final Project for Dibimbing Bootcamp - Golang Back-End development