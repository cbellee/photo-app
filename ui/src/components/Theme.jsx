import { createTheme } from '@mui/material/styles';

export const myTheme = createTheme({
    components: {
        MuiOutlinedInput: {
            styleOverrides: {
                root: {
                    '&.Mui-focused': {
                        borderWidth: 1,
                    }
                }
            }
        }
    }
});