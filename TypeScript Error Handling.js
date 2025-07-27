printUserExperience(4);
printUserExperience(10);
// don't touch below this line
function printUserExperience(level) {
    if (level < 1 || level > 100) {
        throw new Error("Experience level must be between 1 and 100!");
    }
    console.log("The user has ".concat(level, " years of experience!"));
}
