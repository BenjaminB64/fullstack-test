import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import {NextUIProvider} from "@nextui-org/react";

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
      <NextUIProvider>
          <main className="dark text-foreground bg-background h-screen w-screen ">
            <App />
          </main>
      </NextUIProvider>
  </React.StrictMode>
,
)
