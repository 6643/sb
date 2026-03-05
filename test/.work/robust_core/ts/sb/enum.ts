


export enum Color {
    Red = 1, 
    Green = 2, 
    Blue = 3, 
}

export const IsColor = (v: Color): boolean => {
    switch (v) {
    case Color.Red, Color.Green, Color.Blue:
        return true;
    default:
        return false;
    }
}

export const IsColorList = (v: Color[]): boolean => {
    for (const item of v) {
        if (!IsColor(item)) return false;
    }
    return true;
}

