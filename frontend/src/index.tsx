import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';
import "@radix-ui/themes/styles.css";
import { Theme } from '@radix-ui/themes/dist/cjs/components/theme';

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);
root.render(
    <React.StrictMode>
        <Theme>
            <App />
        </Theme>
    </React.StrictMode>
);
