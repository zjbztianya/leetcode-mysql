# [175. Combine Two Tables](https://leetcode-cn.com/problems/combine-two-tables/)
Table: <code>Person</code>


<pre>+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| PersonId    | int     |
| FirstName   | varchar |
| LastName    | varchar |
+-------------+---------+
PersonId is the primary key column for this table.
</pre>

Table: <code>Address</code>


<pre>+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| AddressId   | int     |
| PersonId    | int     |
| City        | varchar |
| State       | varchar |
+-------------+---------+
AddressId is the primary key column for this table.
</pre>



Write a SQL query for a report that provides the following information for each person in the Person table, regardless if there is an address for each of those people:


<pre>FirstName, LastName, City, State
</pre>

## 解题思路## 可能的變化