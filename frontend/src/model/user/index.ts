export const PERMISSION = {
    Player: 0,
    GM: 1,
    Arbiter: 2,
    Admin: 3
} as const

export type PermissionTitle = keyof typeof PERMISSION;

export type PermissionValue = typeof PERMISSION[PermissionTitle]

export type PermissionName = "Игрок" | "ГМ" | "Арбитр" | "Админ"

export const PermissionNameByValue: Record<PermissionValue, PermissionName> = {
    [PERMISSION.Admin]: "Админ",
    [PERMISSION.Arbiter]: "Арбитр",
    [PERMISSION.GM]: "ГМ",
    [PERMISSION.Player]: "Игрок"
}

export const PermissionTitleByValue: Record<PermissionValue, PermissionTitle> = {
    [PERMISSION.Admin]: "Admin",
    [PERMISSION.Arbiter]: "Arbiter",
    [PERMISSION.GM]: "GM",
    [PERMISSION.Player]: "Player"
}

export const PermissionValueByName: Record<PermissionName, PermissionValue> = {
    "Админ": 3,
    "Арбитр": 2,
    "ГМ": 1,
    "Игрок": 0
}

/**
 * Схема хранения данных о пользователе в Cloud Firestore
 */
export type FirestoreUserData = {
    name: string,
    permission: PermissionValue,
    email: string
}

export const PermissionNames: PermissionName[] = ["Игрок", "ГМ", "Арбитр", "Админ"]