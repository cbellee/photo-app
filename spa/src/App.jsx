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

/**
 * If a user is authenticated the ProfileContent component above is rendered. Otherwise a message indicating a user is not authenticated is rendered.
 */

const MainContent = () => {
    const [tabIndex, setTabIndex] = React.useState(0);

    const handleTabChange = (event, newTabIndex) => {
        setTabIndex(newTabIndex);
    };
    return (
        <div className="App">
            <ThemeProvider theme={myTheme}>
                <Routes>
                    <Route path="/collections" element={<CollectionList />} />
                    <Route path="/collections/albums" element={<AlbumList />} />
                </Routes>
                <AuthenticatedTemplate>
                    <Tabs
                        value={tabIndex}
                        onChange={handleTabChange}
                        textColor="secondary"
                        indicatorColor="secondary"
                        aria-label="secondary tabs example"
                    >
                        <Tab label="Collections">
                        </Tab>
                        <Tab label="Upload">
                        </Tab>
                    </Tabs>
                    <Box sx={{ padding: 2 }}>
                        {tabIndex === 0 && (
                            <Box>
                                <CollectionList />
                            </Box>
                        )}
                        {tabIndex === 1 && (
                            <Box>
                                <ImageUpload />
                            </Box>
                        )}
                    </Box>
                </AuthenticatedTemplate>
                <UnauthenticatedTemplate>

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
