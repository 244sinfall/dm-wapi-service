import {PermissionValue} from "../user/types";

export interface Types {
    menuName: string,
    menuRoute?: string,
    accessLevel?: PermissionValue
    action?: () => void
}