


export enum Level {
    Low = 1, 
    Mid = 2, 
    High = 3, 
}

export const IsLevel = (v: Level): boolean => {
    switch (v) {
    case Level.Low, Level.Mid, Level.High:
        return true;
    default:
        return false;
    }
}

export const IsLevelList = (v: Level[]): boolean => {
    for (const item of v) {
        if (!IsLevel(item)) return false;
    }
    return true;
}

