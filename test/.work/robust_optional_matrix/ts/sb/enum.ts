


export enum Kind {
    A = 1, 
    B = 2, 
}

export const IsKind = (v: Kind): boolean => {
    switch (v) {
    case Kind.A, Kind.B:
        return true;
    default:
        return false;
    }
}

export const IsKindList = (v: Kind[]): boolean => {
    for (const item of v) {
        if (!IsKind(item)) return false;
    }
    return true;
}

