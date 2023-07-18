import React from 'react';
import PageLayout from "./components/PageLayout";
import Pages from "./components/Pages";
import "./styles/App.css";
import { MsalProvider } from "@azure/msal-react";
import { useNavigate } from 'react-router-dom';
import { CustomNavigationClient } from "./components/NavigationClient";


function App({ pca }) {
    const history = useNavigate();
    const navigationClient = new CustomNavigationClient(history);
    pca.setNavigationClient(navigationClient);

    return (
        <div className="App">
            <MsalProvider instance={pca}>
                <PageLayout />
                    <Pages />
            </MsalProvider>
        </div>
    );
};

export default App;
