import React, { useState, useEffect } from 'react';
import { AuthenticatedTemplate, UnauthenticatedTemplate, useMsal } from "@azure/msal-react";
import { PageLayout } from "./components/PageLayout";
import { CollectionGrid } from './components/CollectionGrid';
import "./styles/App.css";
import { ThemeProvider } from '@mui/material'
import { myTheme } from './components/Theme';
import PhotoList from './components/PhotoList';

/**
 * If a user is authenticated the ProfileContent component above is rendered. Otherwise a message indicating a user is not authenticated is rendered.
 */
const MainContent = () => {
    return (
        <div className="App">
            <ThemeProvider theme={myTheme}>
                <AuthenticatedTemplate>
                    <PhotoList />
                </AuthenticatedTemplate>
                <UnauthenticatedTemplate>
                    <PhotoList />
                    {/*   <h5 className="card-title">Please sign-in</h5> */}
                </UnauthenticatedTemplate>
            </ThemeProvider>
        </div>
    );
};

export default function App() {
    return (
        <PageLayout>
            <MainContent />
        </PageLayout>
    );
}
