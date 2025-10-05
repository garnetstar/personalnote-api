# React Frontend# React Frontend



This is the React TypeScript frontend for the Simple Go API project.This is the React TypeScript frontend for the Simple Go API project.



## Features## Features



- 🎯 **Interactive API Testing** - Test all API endpoints through a clean web interface- 🎯 **Interactive API Testing** - Test all API endpoints through a clean web interface

- 📱 **Responsive Design** - Works on desktop and mobile devices  - 📱 **Responsive Design** - Works on desktop and mobile devices  

- 🔄 **Real-time Updates** - Live API status and data fetching- 🔄 **Real-time Updates** - Live API status and data fetching

- 🎨 **Modern UI** - Clean, card-based interface with hover effects- 🎨 **Modern UI** - Clean, card-based interface with hover effects

- 🔍 **Article Filtering** - Filter articles by category- 🔍 **Article Filtering** - Filter articles by category

- 👤 **User Management** - Create users through the web interface- 👤 **User Management** - Create users through the web interface



## Available Scripts## Available Scripts



### `npm start`### `npm start`

Runs the app in development mode at [http://localhost:3000](http://localhost:3000).Runs the app in development mode at [http://localhost:3000](http://localhost:3000).



### `npm run build`### `npm run build`

Builds the app for production to the `build` folder.Builds the app for production to the `build` folder.



### `npm test`### `npm run eject`

Launches the test runner in interactive watch mode.

**Note: this is a one-way operation. Once you `eject`, you can’t go back!**

### `npm run eject`

⚠️ **Note: This is a one-way operation. Once you eject, you can't go back!**If you aren’t satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.



## API IntegrationInstead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you’re on your own.



The frontend connects to the Go API running on `http://localhost:8080` and provides:You don’t have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn’t feel obligated to use this feature. However we understand that this tool wouldn’t be useful if you couldn’t customize it when you are ready for it.



- **Health Check**: Test API availability## Learn More

- **Articles List**: Browse all available articles

- **Article Details**: View individual article informationYou can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

- **Article Filtering**: Filter articles by category

- **User Creation**: Add new users to the systemTo learn React, check out the [React documentation](https://reactjs.org/).


## Technologies Used

- **React 18** with TypeScript
- **CSS3** with Flexbox and Grid
- **Fetch API** for HTTP requests
- **Responsive Design** principles

## Development

1. Make sure the Go API is running on port 8080
2. Install dependencies: `npm install`
3. Start development server: `npm start`
4. Open http://localhost:3000 in your browser

## Docker Support

The frontend can be built and deployed using Docker:

```bash
# Build the Docker image
docker build -t simple-go-api-frontend .

# Run the container
docker run -p 3000:3000 simple-go-api-frontend
```

Or use Docker Compose from the project root to run the full stack.