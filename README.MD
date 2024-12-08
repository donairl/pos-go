# Project Overview

This project is a web-based application designed to manage sales transactions, view reports, and monitor inventory. The application is built using HTML, CSS, and JavaScript, and it interacts with a backend API to fetch and manipulate data.

## Features

- **Product Selection**: Users can search and select products to add to their cart.
- **Cart Management**: Users can view the items in their cart, select a payment method, and enter the payment amount.
- **Transaction Processing**: Users can process transactions, and the application will handle payment validation and stock updates.
- **Reports**: Users can view daily and monthly sales summaries, export reports, and view recent transactions.
- **Dashboard**: Users can access quick actions, view today's sales and transactions, and monitor low stock items.

## File Structure

- `views/transactions/index.html`: Handles product selection, cart management, and transaction processing.
- `views/reports/index.html`: Displays sales summaries, allows report export, and shows recent transactions.
- `views/dashboard/index.html`: Provides quick actions, today's summary, and low stock alerts.

## How to Use

1. **Product Selection**:
    - Search for products using the search bar.
    - Click on a product to add it to the cart.

2. **Cart Management**:
    - View the items in the cart.
    - Select a payment method (cash or card).
    - Enter the payment amount.
    - Click "Process Transaction" to complete the purchase.

3. **Reports**:
    - View daily and monthly sales summaries.
    - Export reports by selecting a date range and clicking "Export CSV".
    - View recent transactions and their details.

4. **Dashboard**:
    - Access quick actions to create new transactions, manage products, and view reports.
    - View today's sales and transactions summary.
    - Monitor low stock items and take necessary actions.

## API Endpoints

- `/api/products`: Fetches the list of products.
- `/api/products/search`: Searches for products based on the input query.
- `/api/reports/sales`: Fetches sales data (daily or monthly).
- `/api/reports/export`: Exports sales data as a CSV file.
- `/api/transactions`: Fetches recent transactions and processes new transactions.

## Scripts

- **Product Management**:
    - `addToCart(id, name, price)`: Adds a product to the cart.
    - `calculateChange()`: Calculates the change based on the payment amount.
    - `processTransaction()`: Processes the transaction and updates the backend.
    - `clearCart()`: Clears the cart after a successful transaction.

- **Report Management**:
    - `loadSummary()`: Loads daily and monthly sales summaries.
    - `viewTransaction(id)`: Fetches and displays transaction details.

- **Dashboard Management**:
    - `loadDashboardData()`: Loads today's sales, transactions, and low stock items.
    - `setInterval(loadDashboardData, 5 * 60 * 1000)`: Refreshes dashboard data every 5 minutes.

## Installation

1. Clone the repository.
2. Open the HTML files in a web browser to view the application.
3. Ensure the backend API is running and accessible for data fetching and manipulation.

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes and commit them (`git commit -m 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a new Pull Request.

## License

This project is licensed under the MIT License.

## Acknowledgements

- Thanks to the contributors and the open-source community for their support and contributions.

## Contact

For any questions or inquiries, please contact [your-email@example.com].
