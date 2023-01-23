import React from 'react';
import { AuthenticatedTemplate, UnauthenticatedTemplate } from "@azure/msal-react";
import { PageLayout } from "./components/PageLayout";
import "./styles/App.css";
import { ThemeProvider, Tabs, Tab, Box } from '@mui/material'
import { myTheme } from './components/Theme';
import { ImageUpload } from './components/ImageUpload'
import { Route, Routes } from 'react-router-dom';
import CollectionList from './components/CollectionList';
import AlbumList from './components/AlbumList';
import PhotoList from './components/PhotoList';

const MainContent = () => {
    const [tabIndex, setTabIndex] = React.useState(0);

    const handleTabChange = (event, newTabIndex) => {
        setTabIndex(newTabIndex);
    };
    
    return (
        <div className="App">
            <ThemeProvider theme={myTheme}>
                <AuthenticatedTemplate>
                    <Routes>
                        <Route path="/" element={<CollectionList />} />
                        <Route path="/:collection" element={<AlbumList />} />
                        <Route path="/:collection/:album" element={<PhotoList />} />
                        <Route path="/uploads" element={<ImageUpload />} />
                    </Routes>
                </AuthenticatedTemplate>
                <UnauthenticatedTemplate>
                    <Routes>
                        <Route path="/" element={<CollectionList />}></Route>
                        <Route path="/:collection" element={<AlbumList />} />
                        <Route path="/:collection/:album" element={<PhotoList />} />
                    </Routes>
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
