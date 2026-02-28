import styles from "./Avatar.module.css"

interface AvatarProps {
  name: string
}

function getInitials(name: string): string {
  const pascalCaseMatches = name.match(/[A-Z][a-z]*/g)
  if (pascalCaseMatches && pascalCaseMatches.length >= 2) {
    return pascalCaseMatches[0][0] + pascalCaseMatches[1][0]
  }
  return name[0].toUpperCase()
}

export default function Avatar({ name }: AvatarProps) {
  const initials = getInitials(name)

  return (
    <div className={styles.avatar} role="img" aria-label={`${name} avatar`}>
      <span>{initials}</span>
    </div>
  )
}
